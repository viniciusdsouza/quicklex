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

var gdCmd = &cobra.Command{
	Use:   "gd",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		getAllDefinitions(args[0])
	},
}

func init() {
	rootCmd.AddCommand(gdCmd)
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
	pterm.DefaultBasicText.Println("'" + pterm.LightMagenta(wordDefinitions[0].Word) + "'")
	fmt.Printf("%v\n", wordDefinitions[0].Word)
	fmt.Printf("1. %v - %v\n", wordDefinitions[0].Meanings[0].PartOfSpeech, wordDefinitions[0].Meanings[0].Definitions[0].Definition)
	if wordDefinitions[0].Meanings[0].Definitions[0].Example != "" {
		fmt.Printf(" Ex: %v\n", wordDefinitions[0].Meanings[0].Definitions[0].Example)
	}
	fmt.Printf("2. %v - %v\n", wordDefinitions[0].Meanings[0].PartOfSpeech, wordDefinitions[0].Meanings[1].Definitions[0].Definition)
	if wordDefinitions[0].Meanings[0].Definitions[0].Example != "" {
		fmt.Printf(" Ex: %v\n", wordDefinitions[0].Meanings[1].Definitions[0].Example)
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
