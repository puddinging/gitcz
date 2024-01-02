package cmd

import (
	"embed"
	"encoding/json"
	"log"
	"strings"
)

type CzType struct {
	Type    string `json:"code"`
	Message string `json:"description"`
}

var CzTypeList []CzType

//go:embed data/data.json
var content embed.FS

func Init() {
	data, readErr := content.ReadFile("data/data.json")
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.NewDecoder(strings.NewReader(string(data))).Decode(&CzTypeList)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
