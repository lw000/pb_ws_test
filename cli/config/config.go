package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Cfg configStruct
}

func NewConfig() *Config {
	return &Config{}
}

func LoadConfig(file string) (*Config, error) {
	cfg := &Config{}
	err := cfg.Load(file)
	return cfg, err
}

func (c *Config) Load(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &c.Cfg); err != nil {
		return err
	}

	return nil
}

func (c Config) String() string {
	return fmt.Sprintf("{Count:%d}", c.Cfg.Count)
}
