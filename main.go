package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-wordwrap"
	"github.com/pborman/getopt"
	openai "github.com/spideyz0r/openai-go"
)

const (
	model = "gpt-3.5-turbo"
)

func main() {
	help := getopt.BoolLong("help", 'h', "display this help")
	apiKey := getopt.StringLong("api", 'a', "", "API key (default: OPENAI_API_KEY environment variable)")
	temperature := getopt.StringLong("temperature", 't', "0.8", "temperature (default: 0.8)")
	system_role := getopt.StringLong("system-role", 'r', "You're an expert in everything.", "system role (default: You're an expert in everything. You like speaking.)")
	output_width := getopt.StringLong("output-width", 'w', "80", "output width (default: 80)")
	delim := getopt.StringLong("delimiter", 'd', "\n", "set the delimiter for the user input (default: new line)")
	stdin_input := getopt.BoolLong("stdin", 's', "read the message from stdin and exit (default: false)")

	getopt.Parse()

	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	if *apiKey == "" {
		*apiKey = os.Getenv("OPENAI_API_KEY")
	}

	t, err := strconv.ParseFloat(*temperature, 32)
	if err != nil {
		log.Fatal(err)
	}
	w, err := strconv.Atoi(*output_width)
	if err != nil {
		log.Fatal(err)
	}

	client := openai.NewOpenAIClient(*apiKey)
	messages := []openai.Message{
		{
			Role:    "system",
			Content: *system_role,
		},
	}

	if *stdin_input {
		userInput, err := ioutil.ReadAll(os.Stdin)
		messages = append(messages, openai.Message{
			Role:    "user",
			Content: string(userInput),
		})
		if err != nil {
			log.Fatal(err)
		}
		output, err := sendMessage(client, messages, float32(t), uint(w))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
		os.Exit(0)
	}

	for {
		userInput := getUserInput(*delim)
		userInput = strings.TrimSpace(userInput)
		if userInput == "" {
			continue
		}
		messages = append(messages, openai.Message{
			Role:    "user",
			Content: userInput,
		})

		output, err := sendMessage(client, messages, float32(t), uint(w))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	}
}

func sendMessage(client *openai.OpenAIClient, messages []openai.Message, t float32, w uint) (string, error) {
	completion, err := client.GetCompletion(model, messages, t)
	if err != nil {
		return "", err
	} else {
		return wordwrap.WrapString(fmt.Sprintf("Bot: %v\n", completion.Choices[0].Message.Content), w), nil
	}
}

func getUserInput(delim string) string {
	d := delim
	if delim == "\n" {
		d = "new line"
	}

	fmt.Printf("You (press %s to finish): ", d)
	reader := bufio.NewReader(os.Stdin)
	var input string
	for {
		text, _ := reader.ReadString('\n')
		input += text
		if strings.Contains(input, delim) {
			break
		}
	}
	return input
}
