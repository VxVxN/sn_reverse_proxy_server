package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"reverse_proxy_server/app/tools"
)

type service struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Path    string `yaml:"path"`
}

type config struct {
	Address string `yaml:"address"`

	Services []service `yaml:"services"`

	IsTrace bool `yaml:"trace"`
}

var Cfg *config

func InitConfig(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can't open config: %v", err)
	}
	defer tools.CloseFile(jsonFile)

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("can't read config: %v", err)
	}

	Cfg = &config{}
	if err = yaml.Unmarshal(byteValue, Cfg); err != nil {
		return fmt.Errorf("can't unmarshal config: %v", err)
	}
	return nil
}
