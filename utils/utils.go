package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadConfig() *Config {
	jsonFile, err := os.Open("config.json")

	if err != nil {
		panic(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var cfg Config
	
	json.Unmarshal(byteValue, &cfg)
	defer jsonFile.Close()

	return &cfg
}