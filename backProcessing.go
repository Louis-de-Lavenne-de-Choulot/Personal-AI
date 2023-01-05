package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bregydoc/gtranslate"
)

// var speechFile string = "languagesSpeech"
var translationFile string = "languagesTranslation"

func Serv(input string, lang string, fromLang string) string {
	println("Processing " + input)
	if lang == "" {
		lang = "en"
	}
	if fromLang == "" {
		fromLang = "en"
	}
	println(fromLang)
	//translate input in english
	input = Translate(input, "en", fromLang)
	println(input + "_______________________________________")
	v, fromLang := witAIHandler(input, "")
	println(fromLang)
	if fromLang == "" {
		fromLang = "en"
	}
	return v
}

// func Say(input string, lang string) string {
// 	println("Saying " + input)
// 	//search for language in languagesSpeech.csv
// 	lang = languageSearch(lang, speechFile)
// 	speech := htgotts.Speech{Folder: "audio", Language: lang}
// 	speech.Speak(input)
// 	return input
// }

func Translate(input string, lang string, fromLang string) string {
	//fromLang is optional if not set, it will be set to english
	if fromLang != "en" {
		fromLang = languageSearch(fromLang, translationFile)
	}
	lang = languageSearch(lang, translationFile)
	println("Translating from " + fromLang + " to " + lang)
	translated, err := gtranslate.TranslateWithParams(
		input,
		gtranslate.TranslationParams{
			From: fromLang,
			To:   lang,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return translated
}

func languageSearch(input string, file string) string {
	if len(input) < 3 {
		return input
	}
	//input to lower case
	input = strings.ToLower(input)
	//loop 2 times to search for exact match
	for i := 0; i < 2; i++ {
		//search for language in languages.csv
		csvFile, _ := os.Open(file + ".csv")
		defer csvFile.Close()
		reader := csv.NewReader(bufio.NewReader(csvFile))
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}
			if strings.ToLower(line[i]) == input {
				return line[2]
			}
		}
	}
	return partialLanguageSearch(input, file)
}

func partialLanguageSearch(input string, file string) string {
	//input to lower case
	input = strings.ToLower(input)
	//loop 2 times to search for exact match
	for i := 0; i < 2; i++ {
		//search for language in languages.csv
		csvFile, _ := os.Open(file + ".csv")
		defer csvFile.Close()
		reader := csv.NewReader(bufio.NewReader(csvFile))
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				log.Fatal(error)
			}
			if strings.Contains(input, strings.ToLower(line[i])) {
				return line[2]
			}
		}
	}
	//return language code
	return "en"
}
