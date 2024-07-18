package freeplay

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const DefaultAPIKey = "FREEPLAY_API_KEY"
const DefaultBasePath = "/api/v2/"

var ErrMissingAPIKey = errors.New("missing API key. Set the FREEPLAY_API_KEY environment variable or use the WithAPIKey option")

type Client struct {
	apiKey      string
	apiHost     string
	apiBasePath string

	httpClient *http.Client
	debug      bool
}

type ClientOption func(*Client)

func WithHttpClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

func WithAPIURL(host, basePath string) ClientOption {
	return func(c *Client) {
		c.apiHost = host
		c.apiBasePath = basePath
	}
}

func WithDebug() ClientOption {
	return func(c *Client) {
		c.debug = true
	}
}

func NewClient(baseHost string, opts ...ClientOption) (*Client, error) {
	envKey := os.Getenv(DefaultAPIKey)
	c := &Client{
		httpClient:  http.DefaultClient,
		apiKey:      envKey,
		apiHost:     baseHost,
		apiBasePath: DefaultBasePath,
		debug:       false,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	return c, nil
}

func (c *Client) Debug(format string, args ...any) {
	if c.debug {
		fmt.Printf("[FREEPLAY] %s\n", fmt.Sprintf(format, args...))
	}
}

func (c *Client) AuthGet(url string, body io.Reader) (*http.Response, error) {
	return c.AuthDo("GET", url, body)
}

func (c *Client) AuthPost(url string, body io.Reader) (*http.Response, error) {
	return c.AuthDo("POST", url, body)
}

func (c *Client) AuthDo(verb, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(verb, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	c.Debug("req: %v", req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp, nil
	}

	defer resp.Body.Close()
	result, _ := io.ReadAll(resp.Body)
	return nil, fmt.Errorf("unexpected status code: %d\n\n%s", resp.StatusCode, string(result))
}
