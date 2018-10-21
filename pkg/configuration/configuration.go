package configuration

import (
	"encoding/json"
	"os"
)

type Configuration struct {
		ThisPortNumber         string
		DbURL                  string
		DatabaseName           string
		CollectionName         string
}

func New(path string) (c Configuration, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(&c)
	return
}
