package quasar

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	UUID    string `json:"uuid"`
	Service ServiceConfig
}

type ServiceConfig struct {
	SenderBind   string `json:"sender_bind"`
	ReceiverBind string `json:"receiver_bind"`
}

func NewConfigFromJSONFile(path string) (*Config, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	if err := json.Unmarshal(input, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
