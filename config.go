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
	// SendAddr is the host and port in a TCP URL format: tcp://localhost:61124.
	// If the SendAddr string is blank, the default tenyks port is used.
	SendAddr string `json:"send_addr"`
	// RecvAddr is the host and port in a TCP URL format: tcp://localhost:61124.
	// If the RecvAddr string is blank, the default tenyks port is used.
	RecvAddr string `json:"recv_addr"`
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
