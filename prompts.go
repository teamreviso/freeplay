package freeplay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
)

type Prompt struct {
	FormatVersion           int       `json:"format_version"`
	FormattedContent        []Message `json:"formatted_content"`
	Content                 []Message `json:"content"`
	Metadata                Metadata  `json:"metadata"`
	PromptTemplateID        string    `json:"prompt_template_id"`
	PromptTemplateName      string    `json:"prompt_template_name"`
	PromptTemplateVersionID string    `json:"prompt_template_version_id"`
	SystemContent           *string   `json:"system_content"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type MetadataParams struct {
	ResponseFormat ResponseFormat `json:"response_format"`
	MaxTokens      int            `json:"max_tokens"`
	Temperature    float64        `json:"temperature"`
	TopP           float64        `json:"top_p"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type Metadata struct {
	Flavor       string            `json:"flavor"`
	Model        string            `json:"model"`
	Params       MetadataParams    `json:"params"`
	Provider     string            `json:"provider"`
	ProviderInfo map[string]string `json:"provider_info"`
}

type GetPromptParams struct {
	ProjectID    string
	TemplateName string
	Environment  string
	Formatted    bool

	Data any
}

type GetAllPromptsParams struct {
	ProjectID string
}

func (c *Client) GetAllPrompts(projectID string) ([]Prompt, error) {
	apiURL, err := url.Parse(c.apiHost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %v", err)
	}
	apiURL.Path = path.Join(c.apiBasePath, "projects", projectID, "/prompt-templates/all/latest")

	resp, err := c.AuthGet(apiURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompts: %v", err)
	}
	defer resp.Body.Close()

	var result map[string][]Prompt
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result["prompt_templates"], nil
}

func (c *Client) GetLatestPrompt(
	projectID, templateName string,
	formatted bool,
	data any,
) (*Prompt, error) {
	return c.GetPrompt(projectID, templateName, formatted, "", data)
}

func (c *Client) GetPrompt(
	projectID string,
	templateName string,
	formatted bool,
	environment string,
	data any,
) (*Prompt, error) {
	apiURL, err := url.Parse(c.apiHost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %v", err)
	}
	apiURL.Path = path.Join(c.apiBasePath, "projects", projectID, "prompt-templates", "name", templateName)

	urlParams := url.Values{}
	if environment != "" {
		urlParams.Add("environment", environment)
	}
	urlParams.Add("format", "true")
	apiURL.RawQuery = urlParams.Encode()

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %v", err)
	}

	resp, err := c.AuthPost(apiURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt: %v", err)
	}
	defer resp.Body.Close()

	var result Prompt
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}
