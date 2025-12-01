package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func (c Config) SetUser(name string) {
	c.CurrentUserName = name

	jsonData, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0777); err != nil {
		log.Fatalf("Err: %v", err)
	}

}

func Read() Config {

	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatalf("Err: %v", err)
	}
	file, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatalf("Err: %v", err)
	}

	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Err: %v", err)
	}
	return config
}

func getConfigFilePath() (string, error) {
	usrDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := fmt.Sprintf("%s/%s", usrDir, configFileName)

	return filePath, nil
}
