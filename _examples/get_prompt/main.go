package main

import (
	"fmt"
	"os"

	"github.com/teamreviso/freeplay"
)

func main() {
	client, err := freeplay.NewClient("https://reviso.freeplay.ai")
	if err != nil {
		panic(err)
	}

	prompt, err := client.GetLatestPrompt(
		os.Getenv("FREEPLAY_PROJECT_ID"),
		"drafts convo",
		true,
		map[string]string{},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(prompt)
}
