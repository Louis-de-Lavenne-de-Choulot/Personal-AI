package main

//file for witAI
import (
	"fmt"
	"os"
	"strings"

	"github.com/krognol/go-wolfram"
	witai "github.com/wit-ai/wit-go/v2"
)

var witClient *witai.Client

// function WitAISend to send a message to wit.ai and get the response
func witAIHandler(message string, key string) (string, string) {
	//if key is empty, set it to the environment variable WIT_AI_ACCESS_TOKEN
	if key == "" {
		key = os.Getenv("WIT_AI_ACCESS_TOKEN")
	}
	println("wit.ai key: ", key)
	witClient = witai.NewClient(key)
	resp, err := witClient.Parse(&witai.MessageRequest{Query: message})
	if err != nil {
		println("error: ", err)
		return "Error at client message", ""
	}
	var (
		topIntent           witai.MessageIntent
		topIntentKey        string
		confidenceThreshold float64 = 0.5
	)

	for _, intent := range resp.Intents {
		if intent.Confidence > confidenceThreshold && intent.Confidence > topIntent.Confidence {
			topIntent = intent
			topIntentKey = intent.Name
		}
	}
	//print resp.Text
	println("_____________________________resp.Text: " + resp.Text)

	return reply(topIntentKey, topIntent, message)
}

// function reply to reply to the user
func reply(entityKey string, entity witai.MessageIntent, message string) (string, string) {
	println("_____________________________" + entityKey + "____________________________")
	switch entityKey {
	//switch case for the intents : Calculation,I_am_the_greatest_AI_created_by_Louis_de_Lavenne_de_Choulot,Thank,change_language,greetings,greetings_question,say,search,translate,wit/cancel,wit/stop
	case "Calculation":
		message = strings.ReplaceAll(message, "plus", "+")
		message = strings.ReplaceAll(message, "x", "*")
		message = strings.ReplaceAll(message, "times", "*")
		message = strings.ReplaceAll(message, "time", "*")
		message = strings.ReplaceAll(message, "multiply", "*")
		message = strings.ReplaceAll(message, "multiplied by", "*")
		message = strings.ReplaceAll(message, "divided by", "/")
		message = strings.ReplaceAll(message, "divide", "/")
		message = strings.ReplaceAll(message, "minus", "-")
		message = strings.ReplaceAll(message, "substract", "-")
		message = strings.ReplaceAll(message, "substracted by", "-")
		message = strings.ReplaceAll(message, "to the power of", "^")
		message = strings.ReplaceAll(message, "to the power", "^")
		message = strings.ReplaceAll(message, "power", "^")
		message = strings.ReplaceAll(message, "square root of", "sqrt")
		message = strings.ReplaceAll(message, "square root", "sqrt")
		message = strings.ReplaceAll(message, "square", "^2")
		message = strings.ReplaceAll(message, "cubic", "^3")
		message = strings.ReplaceAll(message, "cubic root of", "cbrt")
		message = strings.ReplaceAll(message, "cubic root", "cbrt")
		//replace all words by empty string
		message = strings.ReplaceAll(message, "calculate", "")
		message = strings.ReplaceAll(message, "'s", "")
		message = strings.ReplaceAll(message, "what", "")
		message = strings.ReplaceAll(message, "does", "")
		message = strings.ReplaceAll(message, "do", "")
		message = strings.ReplaceAll(message, "is", "")
		message = strings.ReplaceAll(message, "the", "")
		message = strings.ReplaceAll(message, "of", "")
		message = strings.ReplaceAll(message, "equal", "")
		message = strings.ReplaceAll(message, "equals", "")
		message = strings.ReplaceAll(message, "to", "")

		//call the function calculation and return the result if no error and the error if error is not nil
		result, err := calculation(message)
		if err != nil {
			return err.Error(), ""
		}
		return fmt.Sprintf("%f", result), ""
	case "I_am_the_greatest_AI_created_by_Louis_de_Lavenne_de_Choulot":
		return "I am the greatest AI created by Louis de Lavenne de Choulot de Chabaud-la-Tour", ""
	case "change_language":
		//Get list of devices
		devices := GetAvailableDevices()
		println("_____________________________devices: ", devices)
		SkipToNext(devices[0].ID)
		return "I can't do that yet", ""
	case "greetings":
		return "Hello", ""
	case "greetings_question":
		return "I'm fine, and you?", ""
	case "Thank":
		return "You're welcome", ""
	case "goodbye":
		return "Bye", ""
	case "say":
		return Say(message, "en"), ""
	case "search":
		return Search(message, os.Getenv("WOLFRAM_APP_ID")), ""
	case "Translation":
		//find last occurence of "to" in message and get the substring after it
		lastOccurence := message[strings.LastIndex(message, "to")+3:]
		return Translate(message, lastOccurence, "en"), lastOccurence
	default:
		return "I don't understand", ""
	}
}

// Search function using wolfram alpha
func Search(message string, key string) string {
	client := &wolfram.Client{AppID: key}
	resp, err := client.GetSpokentAnswerQuery(message, wolfram.Metric, 1000)
	if err != nil {
		println("error: ", err)
		return "Error at client message"
	}
	return resp
}
