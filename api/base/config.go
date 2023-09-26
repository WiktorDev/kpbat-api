package kpbatApi

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Bind     string         `json:"bind"`
	Database DatabaseConfig `json:"database"`
}

type DatabaseConfig struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func throwError() {
	fmt.Println("Can't start app because config.json file not found!")
	os.Exit(0)
}
func LoadConfigFile() Config {
	var config Config
	jsonFile, err := os.Open("config.json")
	if err != nil {
		throwError()
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	if err := json.Unmarshal(byteValue, &config); err != nil {
		throwError()
	}
	return config
}
