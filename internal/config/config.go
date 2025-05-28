package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct { // -export aconfig struct  representing json structure with tags
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// export a "SetUser" method on the "Config" struct that writes the config struct to the  JSON file
// after setting "current_user_name" field
func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := writeGatorConfig(*c)
	if err != nil {
		return err
	}

	return nil
}

// export read function  reads the json file at ~/.gatorconfig.json returns Config struct

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	gatorJson, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}

	var gatorConfig Config

	err = json.Unmarshal(gatorJson, &gatorConfig)
	if err != nil {
		return Config{}, err
	}

	return gatorConfig, nil

}

// helper functions:
func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(homePath, configFileName)
	return fullPath, nil

}

func writeGatorConfig(cfg Config) error {
	gatorJson, err := json.Marshal(cfg)

	if err != nil {
		return err
	}

	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(fullPath, gatorJson, 0666)
	if err != nil {
		return err
	}

	return nil
}
