package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/pborman/getopt"
	"github.com/sap-nocops/duckduckgogo/client"
)

func main() {
	help := getopt.BoolLong("help", 'h', "display this help")
	apiKey := getopt.StringLong("api", 'a', "", "API key (default: OPENAI_API_KEY environment variable)")
	temperature := getopt.StringLong("temperature", 't', "0.8", "temperature (default: 0.8)")
	model := getopt.StringLong("model", 'm', "gpt-4", "gpt chat model (default: gpt-4)")
	system_role := getopt.StringLong("system-role", 'S', "You're an expert in everything.", "system role (default: You're an expert in everything. You like speaking.)")
	output_width := getopt.StringLong("output-width", 'w', "80", "output width (default: 80)")
	delim := getopt.StringLong("delimiter", 'd', "\n", "set the delimiter for the user input (default: new line)")
	stdin_input := getopt.BoolLong("stdin", 's', "read the message from stdin and exit (default: false)")
	internet_access := getopt.BoolLong("internet", 'i', "allow internet access (default: false)")
	debug := getopt.BoolLong("debug", 'D', "debug mode (default: false)")

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

	client := openai.NewClient(option.WithAPIKey(*apiKey))
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(*system_role),
	}

	if *stdin_input {
		userInput, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		message := buildMessage(string(userInput), *apiKey, float32(t), *model, *internet_access, *debug)
		messages = append(messages, openai.UserMessage(string(message)))

		output, err := sendMessage(&client, messages, float32(t), *model)
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

		stop := make(chan bool)
		go spinner(10*time.Millisecond, stop)

		message := buildMessage(userInput, *apiKey, float32(t), *model, *internet_access, *debug)
		messages = append(messages, openai.UserMessage(message))
		output, err := sendMessage(&client, messages, float32(t), *model)
		if err != nil {
			log.Fatal(err)
		}
		stop <- true
		fmt.Println(wordwrap.WrapString(fmt.Sprintf("\nBot: %v\n", output), uint(w)))
	}
}

func sendMessage(client *openai.Client, messages []openai.ChatCompletionMessageParamUnion, t float32, model string) (string, error) {
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       model,
		Temperature: openai.Float(float64(t)),
	})
	if err != nil {
		return "", err
	}
	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from chat completion")
	}

	return chatCompletion.Choices[0].Message.Content, nil
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
	return strings.TrimSuffix(input, delim)
}

func spinner(delay time.Duration, stop chan bool) {
	for {
		for _, r := range `-\|/` {
			select {
			case <-stop:
				fmt.Print("\r")
				return
			default:
				fmt.Printf("\r%c", r)
				time.Sleep(delay)
			}
		}
	}
}

func internetSearch(query string) (string, error) {
	ddg := client.NewDuckDuckGoSearchClient()
	res, err := ddg.SearchLimited(query, 5)

	var concatenatedString string
	for _, r := range res {
		concatenatedString += fmt.Sprintf("Title: %s\nSnippet: %s\n\n",
			r.Title, r.Snippet)
	}
	return concatenatedString, err
}

func isRealtimeQuestion(message, apiKey string, t float32, today, model string) (bool, string) {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	todaydate := time.Now().Format("2006-01-02")
	messageTemplate := ` I need your answer to be in a json format. {\"real-time\": \"boolean\", \"message\": \"message\"}. Don't say anything other than the json, nothing.
I am going to ask you a question, consider that today is %s If this question requires real-time access to data, you will answer in the json format with real-time as true and a short message. You don't have the ability to
access real-time data, so keep that in mind. Question related to a time after your last update will be answered with real-time as true. If it doesn't require real-time, the real-time field will be false and in the msg field you can put the full answer. Please don't answer anything other than the json.
Example: Is the Formula 1 GP today?. After this message I'll send the first question. Just answer with the json.`

	msg_content := fmt.Sprintf(messageTemplate, todaydate)
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(msg_content),
	}
	messages = append(messages, openai.UserMessage(message))
	output, err := sendMessage(&client, messages, float32(t), model)
	if err != nil {
		log.Fatal(err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(fmt.Sprintf("%v", output)), &response); err != nil {
		log.Fatal(err)
	}
	realTime, _ := response["real-time"].(bool)
	respMessage, _ := response["message"].(string)

	return realTime, respMessage
}

func buildMessage(userInput, apiKey string, t float32, model string, internet_access, debug bool) string {
	today_date := time.Now().Format("2006-01-02")
	real_time := false
	msg := ""

	if internet_access {
		real_time, msg = isRealtimeQuestion(userInput, apiKey, float32(t), today_date, model)
		if debug {
			fmt.Printf("Real time: %v\nMessage: %s\n", real_time, msg)
		}
	}

	if !real_time {
		return userInput
	}

	if debug {
		fmt.Printf("Real time mode active for question: %s\n", msg)
	}
	results, err := internetSearch(userInput)

	if err != nil {
		fmt.Println("Error making internet search.")
		log.Fatal(err)
	}
	m := "Today is %s, you don't have the information you need to answer this question. So I made a google search for you. Here are the results:\n\n%s.\n\n Now consider this information and answer the question, but pretend we didn't talk about the Internet search. Just answer the question. Today is %s and this is the original question: %s"
	return fmt.Sprintf(m, today_date, results, today_date, userInput)
}
