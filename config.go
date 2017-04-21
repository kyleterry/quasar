package quasar

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Config holds configuration values passed into quasar on bootstrap.
// Recommended way to use this object is to unmarshal some JSON into an instance.
// You can also build this by hand using flags.
type Config struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	UUID    string        `json:"uuid"`
	Service ServiceConfig `json:"service"`
}

// ServiceConfig holds the addresses for pubsub actions
type ServiceConfig struct {
	// SendAddr is the host and port in a TCP URL format: tcp://localhost:61124.
	// If the SendAddr string is blank, the default tenyks port is used.
	SendAddr string `json:"send_addr"`
	// RecvAddr is the host and port in a TCP URL format: tcp://localhost:61124.
	// If the RecvAddr string is blank, the default tenyks port is used.
	RecvAddr string `json:"recv_addr"`
}

// NewConfigFromJSONFile takes a file path to a JSON file and unmarshals
// it into a new Config instance.
func NewConfigFromJSONFile(path string) (*Config, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file")
	}

	conf := &Config{}
	if err := json.Unmarshal(input, conf); err != nil {
		return nil, errors.Wrap(err, "failed to parse JSON bytes")
	}

	return conf, nil
}
