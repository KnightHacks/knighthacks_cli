package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Auth struct {
		Tokens struct {
			Refresh string `yaml:"refresh"`
			Access  string `yaml:"access"`
		} `yaml:"tokens"`
	} `yaml:"auth"`
}

func (c *Config) Load(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, c)
}

func (c *Config) Save(path string) error {
	marshal, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, marshal, 0644)
}
