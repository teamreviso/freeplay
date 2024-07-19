package main

import (
	"fmt"
	"os"
	"time"

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

	start := time.Now()
	prompt, err := client.GetLatestPrompt(
		os.Getenv("FREEPLAY_PROJECT_ID"),
		"drafts convo",
		true,
		map[string]string{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("took %s\n", time.Since(start))
	fmt.Println(prompt.Content)
}
