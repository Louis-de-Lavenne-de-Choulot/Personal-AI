package main

//file for witAI
import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/krognol/go-wolfram"
	witai "github.com/wit-ai/wit-go/v2"
	"github.com/zmb3/spotify/v2"
)

var (
	witClient *witai.Client
	vol       int = 50
	device    spotify.ID
)

// function WitAISend to send a message to wit.ai and get the response
func witAIHandler(message string, key string) (string, string) {
	//if key is empty, set it to the environment variable WIT_AI_ACCESS_TOKEN
	if key == "" {
		key = os.Getenv("WIT_AI_ACCESS_TOKEN")
	}
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
	message = strings.ReplaceAll(message, "Franck", "")
	switch entityKey {
	case "max_volume":
		searchForDevice()
		if device != "" {
			vol = 100
			SetPlaybackVolume(vol, device)
			return "", ""
		}
		return "", ""
	case "min_volume":
		searchForDevice()
		if device != "" {
			vol = 5
			SetPlaybackVolume(vol, device)
			return "", ""
		}
		return "", ""
	case "wit$decrease_volume":
		searchForDevice()
		//decrease device volume
		if device != "" {
			vol -= 5
			SetPlaybackVolume(vol, device)
			return "", ""
		}
		return "the device was not given", ""
	case "wit$increase_volume":
		searchForDevice()
		//increase device volume
		if device != "" {
			vol += 5
			SetPlaybackVolume(vol, device)
			return "", ""
		}
		return "the device was not given", ""
	case "wit$skip_track":
		searchForDevice()
		if device != "" {
			SkipToNext(device)
			return "", ""
		}
		return "the device was not given", ""
	case "wit$previous_track":
		searchForDevice()
		if device != "" {
			SkipToPrevious(device)
			return "", ""
		}
		return "the device was not given", ""
	case "change_language":
		//change language
		return "sorry, I can't do that yet", ""
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

	case "I_am_the_greatest_AI_created_by_Louis_de_Lavenne_de_Choulot":
		return "I am the greatest AI created by Louis de Lavenne de Choulot de Chabaud-la-Tour", ""
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

func searchForDevice() {
	if device == "" {
		//Get list of devices
		devices := GetAvailableDevices()
		listdev := [][]string{}
		Say("Which device ?", "en")
		//loop through devices and say the name
		for _, rdevice := range devices {
			Say(string(rdevice.Name), "en")
			listdev = append(listdev, []string{string(rdevice.ID), string(rdevice.Name)})
		}
		//scan for the device name
		println("Write device name")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		//for device in listdev contains scanner.Text() set device to device[x]
		for _, dev := range listdev {
			if strings.Contains(dev[1], scanner.Text()) {
				device = spotify.ID(dev[0])
			}
		}
		//wait 10 seconds for the user to answer
		// time.Sleep(10 * time.Second)
	}
}
