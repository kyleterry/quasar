package quasar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexMatcherCanMatch(t *testing.T) {
	matcher := NewRegexMatcher("^(?i)(hi|hello|sup|hey), I'm (?P<name>(.*))$")
	msg := Message{Payload: "hey, I'm Kyle"}

	match := matcher.Match(msg)
	assert.NotNil(t, match)

	assert.Contains(t, match, "name")
	assert.Equal(t, match["name"], "Kyle")
}

func TestRegexMatcherWillReturnNilForNoMatch(t *testing.T) {
	matcher := NewRegexMatcher("^(?i)(hi|hello|sup|hey), I'm (?P<name>(.*))$")
	msg := Message{Payload: "nonsense"}

	match := matcher.Match(msg)

	assert.Nil(t, match)
}

func TestRegexMatcherCanFindOneOutOfMany(t *testing.T) {
	matcher := NewRegexMatcher(
		"^(?i)hi, I'm (?P<name>(.*))$",
		"^testing (?P<thing>(.*))$")
	msg := Message{Payload: "testing something"}

	match := matcher.Match(msg)
	assert.NotNil(t, match)

	assert.Contains(t, match, "thing")
	assert.Equal(t, match["thing"], "something")
}
