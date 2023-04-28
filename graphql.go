package graphql

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type GraphQLClient struct {
	Client   *http.Client
	Endpoint string
	Header   http.Header
}

func NewGraphQLClient(endpoint string, httpClient *http.Client, header http.Header) *GraphQLClient {
	return &GraphQLClient{
		Client:   httpClient,
		Endpoint: endpoint,
		Header:   header,
	}
}

func (c *GraphQLClient) Query(graphqlRequest GraphQLRequest) (*GraphQLResponse, error) {
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
	resp, err := c.Client.Do(req)
	if err != nil {
		log.Printf("error making request: %v", err)
		return nil, err
	}

	// Decode the response body
	var result GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("error decoding response: %v", err)
		return nil, err
	}

	return &result, nil
}

func (c *GraphQLClient) Mutation(graphqlRequest GraphQLRequest) (*GraphQLResponse, error) {
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
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Decode the response body
	var result GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type GraphQLResponse struct {
	Data       any                    `json:"data"`
	Errors     []GraphQLError         `json:"errors"`
	Extensions map[string]interface{} `json:"extensions"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLErrorLocation `json:"locations"`
	Path      string                 `json:"path"`
}

type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}
