package freeplay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
)

type SessionInfo struct {
	CustomMetatdata map[string]string `json:"custom_metadata"`
}

type PromptInfo struct {
	PromptTemplateID        string `json:"prompt_template_id"`
	PromptTemplateVersionID string `json:"prompt_template_version_id"`
	PromptTemplateName      string `json:"prompt_template_name"`
	Environment             string `json:"environment"`
	ModelParameters         struct {
		ResponseFormat string  `json:"response_format"`
		Temperature    float64 `json:"temperature"`
		TopP           float64 `json:"top_p"`
	}
	ProviderInfo map[string]string `json:"provider_info"`
	Provider     string            `json:"provider"`
	Model        string            `json:"model"`
	FlavorName   string            `json:"flavor_name"`
	ProjectID    string            `json:"project_id"`
}

type CallInfo struct {
	StartTime     float64           `json:"start_time"`
	EndTime       float64           `json:"end_time"`
	Model         string            `json:"model"`
	Provider      string            `json:"provider"`
	ProviderInfo  map[string]string `json:"provider_info"`
	LlmParameters map[string]string `json:"llm_parameters"`
}

type TestRunInfo struct {
	TestRunId  string `json:"test_run_id"`
	TestCaseId string `json:"test_case_id"`
}

type CompletionPayload struct {
	Messages    []Message         `json:"messages"`
	Inputs      map[string]string `json:"inputs"`
	SessionInfo *SessionInfo      `json:"session_info,omitempty"`
	PromptInfo  PromptInfo        `json:"prompt_info"`
	CallInfo    *CallInfo         `json:"call_info,omitempty"`
	TestRunInfo *TestRunInfo      `json:"test_run_info,omitempty"`
}

type CompletionResponse struct {
	CompletionID string `json:"completion_id"`
}

func (c *Client) RecordCompletion(
	projectID string,
	sessionID string,
	payload *CompletionPayload,
) (*CompletionResponse, error) {
	apiURL, err := url.Parse(c.apiHost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse API URL: %v", err)
	}
	apiURL.Path = path.Join(c.apiBasePath, "projects", projectID, "sessions", sessionID, "completions")

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %v", err)
	}

	c.Debug("POST %s\n", apiURL.String())
	c.Debug("body:\n%s\n", string(body))

	resp, err := c.AuthPost(apiURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to record completion: %v", err)
	}
	defer resp.Body.Close()

	var result CompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

type TracePayload struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func (c *Client) RecordTrace(
	projectID string,
	sessionID string,
	traceID string,
	payload *TracePayload,
) error {
	apiURL, err := url.Parse(c.apiHost)
	if err != nil {
		return fmt.Errorf("failed to parse API URL: %v", err)
	}
	apiURL.Path = path.Join(c.apiBasePath, "projects", projectID, "sessions", sessionID, "traces/id", traceID)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	c.Debug("POST %s\n", apiURL.String())
	c.Debug("body:\n%s\n", string(body))

	_, err = c.AuthPost(apiURL.String(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to record completion: %v", err)
	}

	return nil
}
