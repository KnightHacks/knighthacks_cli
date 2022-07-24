package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	Auth struct {
		Tokens struct {
			Refresh string `json:"refresh"`
			Access  string `json:"access"`
		} `json:"tokens"`
		UserID string `json:"user_id"`
	} `json:"auth"`
}

func (c *Config) Load(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, c)
}

func (c *Config) Save(path string) error {
	marshal, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, marshal, 0644)
}
