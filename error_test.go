package quasar

import (
	"errors"
	"testing"
)

func TestIsNoMatch(t *testing.T) {
	err := ErrNoMatch

	if !IsNoMatch(err) {
		t.Error("expected: true, got: false")
	}

	err = errors.New("test")

	if IsNoMatch(err) {
		t.Error("expected: false, got: true")
	}
}
