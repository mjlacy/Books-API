package configuration

import (
	"os"
	"encoding/json"
	"fmt"
)

type Configuration struct {
		ThisPortNumber         string
		DbURL                  string
		DatabaseName           string
		CollectionName         string
}

func New(path string) Configuration {
	conf := Configuration{}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening properties file: ", err)
	}
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		fmt.Println("Error decoding properties file ", err)
	}
	return conf
}