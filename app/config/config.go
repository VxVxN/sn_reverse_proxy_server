package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"reverse_proxy_server/app/tools"
)

type config struct {
	WebServerHostname string `json:"WEB_SERVER_HOSTNAME"`
	WebServerPort     int    `json:"WEB_SERVER_PORT"`

	AJAXServerHostname string `json:"AJAX_SERVER_HOSTNAME"`
	AJAXServerPort     int    `json:"AJAX_SERVER_PORT"`

	ReverseProxyServerHostname string `json:"REVERSE_PROXY_HOSTNAME"`
	ReverseProxyServerPort     int    `json:"REVERSE_PROXY_SERVER"`
}

var Cfg *config

func InitConfig() error {
	jsonFile, err := os.Open( "config/main.json")
	if err != nil {
		return fmt.Errorf("can't open config: %v", err)
	}
	defer tools.CloseFile(jsonFile)

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("can't read config: %v", err)
	}

	Cfg = &config{}
	if err = json.Unmarshal(byteValue, Cfg); err != nil {
		return fmt.Errorf("can't unmarshal config: %v", err)
	}
	return nil
}

