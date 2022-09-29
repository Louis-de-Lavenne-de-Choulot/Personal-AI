package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/bregydoc/gtranslate"
	htgotts "github.com/hegedustibor/htgo-tts"
)

var speechFile string = "languagesSpeech"
var translationFile string = "languagesTranslation"

func Serv(mod string, input string, lang string, fromLang string) string {
	if lang == "" {
		lang = "en"
	}
	if fromLang == "" {
		fromLang = "en"
	}
	if strings.Contains(mod, "Say") {
		//translate input in the language
		input = Translate(input, lang, fromLang)
		return Say(input, lang, fromLang)
	} else if strings.Contains(mod, "Discuss") {
		//translate input in english
		input = Translate(input, "en", fromLang)
		return Translate(witAIHandler(input, ""), lang, "en")
	}
	return ("Invalid input")
}

func Say(input string, lang string, fromLang string) string {
	//search for language in languagesSpeech.csv
	lang = languageSearch(lang, false, speechFile)
	speech := htgotts.Speech{Folder: "audio", Language: lang}
	speech.Speak(input)
	return input
}

func Translate(input string, lang string, fromLang string) string {
	//fromLang is optional if not set, it will be set to english
	if fromLang != "en" {
		fromLang = languageSearch(fromLang, true, translationFile)
	}
	//language search if lang length is higher than 3
	if len(lang) > 3 {
		lang = languageSearch(lang, true, translationFile)
	}
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

func languageSearch(input string, fileType bool, file string) string {
	if len(input) < 3 {
		return input
	}
	//input to lower case
	input = strings.ToLower(input)
	var nbre int = 0
	//check if file is has 2 or 3 columns
	if fileType {
		nbre = 1
	}
	//loop 2 times to search for exact match
	for i := 0; i < 1+nbre; i++ {
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
				return line[1+nbre]
			}
		}
	}
	return partialLanguageSearch(input, fileType, file)
}

func partialLanguageSearch(input string, fileType bool, file string) string {
	//input to lower case
	input = strings.ToLower(input)
	var nbre int = 0
	//check if file is has 2 or 3 columns
	if fileType {
		nbre = 1
	}
	//loop 2 times to search for exact match
	for i := 0; i < 1+nbre; i++ {
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
			if strings.Contains(strings.ToLower(line[i]), input) {
				return line[1+nbre]
			}
		}
	}
	//return language code
	return "en"
}
