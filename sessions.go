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
	PromptTemplateVersionID string `json:"prompt_template_version_id"`
	Environment             string `json:"environment"`
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

// https://dev.freeplay.ai/api/v2/projects/8f93dd00-2eb5-4ba2-9354-86d5c6831dfd/sessions/f503c15e-2f0f-4ce4-b443-4c87d0b6435d/completions
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
