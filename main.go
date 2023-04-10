package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	openai "github.com/spideyz0r/openai-go"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewOpenAIClient(apiKey)
	temperature := float32(0.8)
	model := "gpt-3.5-turbo"

	messages := []openai.Message{
		{
			Role: "system",
			Content: "You're an expert in anything.",
		},
	}

	for {
		fmt.Print("You: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		userInput = strings.TrimSpace(userInput)

		if userInput == "" {
			continue
		}
		messages = append(messages, openai.Message{
			Role: "user",
			Content: userInput,
		})

		completion, err := client.GetCompletion(model, messages, temperature)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Bot: %v\n", completion.Choices[0].Message.Content)
		}
	}
}
