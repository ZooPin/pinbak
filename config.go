package pinbak

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Repository map[string]string `json:"Repository"`
	path       string            `json:"-"`
	configPath string            `json:"-"`
}

const configName = "config"

func LoadConfig(path string) (Config, error) {
	var config Config
	return config.Load(path)
}

func (c *Config) SetPath(path string) {
	c.path = path
	c.configPath = fmt.Sprint(path, "/", configName)
}

func (c Config) Load(path string) (Config, error) {
	var conf Config
	conf.SetPath(path)

	file, err := os.Open(conf.configPath)
	defer file.Close()

	if err != nil {
		return Config{}, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}

func (c Config) Save() error {
	file, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.configPath, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) AddRepository(name string, url string) error {
	if c.Repository == nil {
		c.Repository = make(map[string]string)
	}
	if val, ok := c.Repository[name]; ok {
		return errors.New(fmt.Sprint("Name already exist with value: ", val))
	}
	c.Repository[name] = url

	err := c.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c Config) CheckRepository(name string) bool {
	_, ok := c.Repository[name]
	return ok
}
