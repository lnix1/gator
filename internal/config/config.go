package config

import (
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbUrl 			string `json:"db_url"`
	CurrentUserName 	string `json:"current_user_name"`
}

func (c Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	write(c)
	return nil
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := homePath + configFileName
	return fullPath, nil
}

func write(cfg Config) error {
	writePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(writePath, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("Config file updated: ")
	return nil
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var currentConfig Config
	err = json.Unmarshal(configBytes, &currentConfig)
	if err != nil {
		return Config{}, err
	}

	return currentConfig, nil
}
