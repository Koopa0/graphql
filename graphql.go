package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Client struct {
	HttpClient *http.Client
	Endpoint   string
	Header     http.Header
}

func NewGraphQLClient(endpoint string, httpClient *http.Client, header http.Header) *Client {
	return &Client{
		HttpClient: httpClient,
		Endpoint:   endpoint,
		Header:     header,
	}
}

func (c *Client) Query(graphqlRequest Request) (*Response, error) {
	// Build the request body
	body, _ := json.Marshal(graphqlRequest)

	// Make the request
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Println("error creating request:", err)
		return nil, err
	}

	req.Header = c.Header
	req.Header.Set("Content-Type", "application/json")

	// Get the response
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("error making request: %v", err)
		return nil, err
	}

	// Decode the response body
	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("error decoding response: %v", err)
		return nil, err
	}

	return &result, nil
}

func (c *Client) QueryWithContext(ctx context.Context, graphqlRequest Request, cookies ...*http.Cookie) (*Response, error) {
	// Build the request body
	body, _ := json.Marshal(graphqlRequest)

	// Make the request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Println("error creating request:", err)
		return nil, err
	}

	req.Header = c.Header
	req.Header.Set("Content-Type", "application/json")

	if len(cookies) > 0 {
		for _, v := range cookies {
			req.AddCookie(v)
		}
	}

	// Get the response
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("error making request: %v", err)
		return nil, err
	}

	// Decode the response body
	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("error decoding response: %v", err)
		return nil, err
	}

	return &result, nil
}

func (c *Client) Mutation(graphqlRequest Request) (*Response, error) {
	// Build the request body
	body, _ := json.Marshal(graphqlRequest)

	// Make the request
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = c.Header
	req.Header.Set("Content-Type", "application/json")

	// Get the response
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Decode the response body
	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type Response struct {
	Data       any            `json:"data"`
	Errors     []Error        `json:"errors"`
	Extensions map[string]any `json:"extensions"`
}

type Request struct {
	Query     string         `json:"query,omitempty"`
	Variables map[string]any `json:"variables,omitempty"`
}

type Error struct {
	Message   string          `json:"message"`
	Locations []ErrorLocation `json:"locations"`
	Path      []any           `json:"path"`
}

type ErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
