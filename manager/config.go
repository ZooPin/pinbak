package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Name       string            `json:"name"`
	Email      string            `json:"email"`
	Repository map[string]string `json:"Repository"`
	path       string            `json:"-"`
	configPath string            `json:"-"`
}

const configName = "config.json"

func LoadConfig(path string) (Config, error) {
	var config Config
	return config.Load(path)
}

func (c *Config) SetPath(lPath string) {
	c.path = lPath
	c.configPath = path.Join(lPath, configName)
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
	return conf, err
}

func (c Config) Save() error {
	file, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.configPath, file, 0644)
	return err
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
	return err
}

func (c *Config) RemoveRepository(name string) error {
	if _, ok := c.Repository[name]; !ok {
		return errors.New("Repository not found.")
	}
	delete(c.Repository, name)
	err := c.Save()
	if err != nil {
		return err
	}
	err = os.RemoveAll(path.Join(c.path, name))
	return err
}

func (c Config) CheckRepository(name string) bool {
	_, ok := c.Repository[name]
	return ok
}
