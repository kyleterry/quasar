package quasar

import "errors"

var ErrNoMatch = errors.New("message payload did not match")

func IsNoMatch(err error) bool {
	if err == nil {
		return false
	}
	return err == ErrNoMatch
}
