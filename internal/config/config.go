package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func (c Config) SetUser(name string) error {
	c.CurrentUserName = name

	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, jsonData, 0777); err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.ReadFile(filePath)

	if err != nil {
		return Config{}, err
	}

	var config Config

	if err := json.Unmarshal(file, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	usrDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := fmt.Sprintf("%s/%s", usrDir, configFileName)

	return filePath, nil
}
