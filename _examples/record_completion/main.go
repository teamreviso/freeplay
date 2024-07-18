package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/teamreviso/freeplay"
)

func main() {
	client, err := freeplay.NewClient(
		"https://reviso.freeplay.ai",
		freeplay.WithDebug(),
	)
	if err != nil {
		panic(err)
	}

	completion, err := client.RecordCompletion(
		os.Getenv("FREEPLAY_PROJECT_ID"),
		uuid.NewString(),
		&freeplay.CompletionPayload{
			Messages: []freeplay.Message{
				{
					Content: "hello",
					Role:    "user",
				},
			},
			CallInfo: &freeplay.CallInfo{
				StartTime:     float64(time.Now().Unix()),
				EndTime:       float64(time.Now().Unix()),
				Model:         "gpt-3.5-turbo",
				Provider:      "openai",
				ProviderInfo:  map[string]string{},
				LlmParameters: map[string]string{},
			},
			Inputs: map[string]string{},
			PromptInfo: freeplay.PromptInfo{
				PromptTemplateVersionID: "cabd1a71-d982-48c5-b779-df9a3d1a8380",
				Environment:             "latest",
			},
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(completion)
}
