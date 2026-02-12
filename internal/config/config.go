package config

import (
	"encoding/json"
	"os"
)


type Config struct {
	DBUrl    string `json:"db_url"`
	UserName string `json:"current_user_name"`
}

const configFileName = "/.gatorconfig.json"

func Read() (Config, error) {
	c := Config{}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return c, err
	}

	file, err := os.ReadFile(homeDir + configFileName)
	if err != nil {
		return c, err
	}

	json.Unmarshal(file, &c)
	return c, nil
}

func (c Config) SetUser(u string) error {
	c.UserName = u
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	_, err = os.Create(homeDir + configFileName)
	if err != nil {
		return err
	}

	jsonbytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := os.WriteFile(homeDir + configFileName, jsonbytes, os.ModeAppend); err != nil {
		return err
	}
	return nil
}
