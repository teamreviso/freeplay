package main

import (
	"fmt"
	"os"

	"github.com/teamreviso/freeplay"
)

func main() {
	client, err := freeplay.NewClient(os.Getenv("FREEPLAY_API_HOST"))
	if err != nil {
		panic(err)
	}

	prompts, err := client.GetAllPrompts(os.Getenv("FREEPLAY_PROJECT_ID"))
	if err != nil {
		panic(err)
	}

	for _, prompt := range prompts {
		fmt.Println(prompt)
	}
}
