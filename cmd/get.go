/*
Copyright © 2023 Vinícius Duarte <vduartesantiago@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve word definitions",
	Long: `	The "get" command is a command-line interface (CLI) command that allows users 
	to retrieve the definitions of a given word directly from their terminal 
	or command prompt. It provides a quick and convenient way to access the meaning, 
	definitions, and sometimes additional information about a word without the need 
	to open a browser or search through extensive dictionaries.

	Syntax:
	ql get [word]

	Parameters:
	[word]: The word for which you want to retrieve the definitions. 
	This parameter is required and should be replaced with the actual word you want to look up.`,
	Run: func(cmd *cobra.Command, args []string) {
		getAllDefinitions(args[0])
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}

type WordDefinitions []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio,omitempty"`
	} `json:"phonetics"`
	Origin   string `json:"origin"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example"`
			Synonyms   []any  `json:"synonyms"`
			Antonyms   []any  `json:"antonyms"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func getAllDefinitions(word string) {
	url := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%v", word)
	responseBytes := getDefinitionData(url)
	wordDefinitions := WordDefinitions{}
	if err := json.Unmarshal(responseBytes, &wordDefinitions); err != nil {
		fmt.Println("Could not find this word!")
		return
	}
	pterm.DefaultCenter.Println("\n'" + pterm.LightMagenta(wordDefinitions[0].Word) + "'")
	pterm.DefaultCenter.Println("1. " + wordDefinitions[0].Meanings[0].PartOfSpeech + " - " + wordDefinitions[0].Meanings[0].Definitions[0].Definition)
	if wordDefinitions[0].Meanings[0].Definitions[0].Example != "" {
		pterm.DefaultCenter.Println(" Ex: " + wordDefinitions[0].Meanings[0].Definitions[0].Example + "\n")
	}
	if 1 < len(wordDefinitions[0].Meanings) {
		pterm.DefaultCenter.Println("\n2. " + wordDefinitions[0].Meanings[1].PartOfSpeech + " - " + wordDefinitions[0].Meanings[1].Definitions[0].Definition)
		if wordDefinitions[0].Meanings[1].Definitions[0].Example != "" {
			pterm.DefaultCenter.Println(" Ex: " + wordDefinitions[0].Meanings[1].Definitions[0].Example + "\n")
		}
		return
	}
}

func getDefinitionData(baseAPI string) []byte {
	req, _ := http.NewRequest("GET", baseAPI, nil)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Could not make a request - %v", err)
	}

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Could not read response body - %v", err)
	}

	return responseBytes
}
