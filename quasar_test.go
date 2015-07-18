package quasar

import (
	"testing"
)

func TestNewConnection(t *testing.T) {
	conf := Config{}
	conn := NewConnection(conf)
}
