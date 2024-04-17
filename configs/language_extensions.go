package configs

import (
	_ "embed"
	"encoding/json"
)

//go:embed language_extensions.json
var languagesJson string

type Language struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

func GetLanguages() ([]Language, error) {
	var allLanguages []Language
	err := json.Unmarshal([]byte(languagesJson), &allLanguages)
	if err != nil {
		return nil, err
	}
	return allLanguages, nil
}
