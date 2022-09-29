package main

//file for witAI
import (
	"os"

	"github.com/christianrondeau/go-wit"
	"github.com/krognol/go-wolfram"
)

var witClient *wit.Client
var wolframClient *wolfram.Client

// function WitAISend to send a message to wit.ai and get the response
func witAIHandler(message string, key string) string {
	//if key is empty, set it to the environment variable WIT_AI_ACCESS_TOKEN
	if key == "" {
		key = os.Getenv("WIT_AI_ACCESS_TOKEN")
	}
	witClient = wit.NewClient(key)
	resp, err := witClient.Message(message)
	if err != nil {
		println("error: ", err)
		return "Error at client message"
	}
	var (
		topEntity           wit.MessageEntity
		topEntityKey        string
		confidenceThreshold float64 = 0.5
	)

	for entityKey, entityList := range resp.Entities {
		for _, entity := range entityList {
			if entity.Confidence > confidenceThreshold && entity.Confidence > topEntity.Confidence {
				topEntity = entity
				topEntityKey = entityKey
			}
		}
	}
	return reply(topEntityKey, topEntity, message)
}

// function reply to reply to the user
func reply(entityKey string, entity wit.MessageEntity, message string) string {
	println("entityKey: ", entityKey)
	switch entityKey {
	case "greetings":
		return "Hello"
	case "You're welcome":
		return "You're welcome"
	case "bye":
		return "Bye"
	case "intent":
		return Search(message, os.Getenv("WOLFRAM_APP_ID"))
	default:
		return "I don't understand"
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
