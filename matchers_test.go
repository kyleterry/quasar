package quasar

import "testing"

func TestRegexMatcherCanMatch(t *testing.T) {
	matcher := NewRegexMatcher("^(?i)(hi|hello|sup|hey), I'm (?P<name>(.*))$")
	msg := Message{Payload: "hey, I'm Kyle"}

	match, err := matcher.Match(msg)
	if err != nil {
		t.Error("expected: nil, got: ", err)
	}

	if name, ok := match["name"]; !ok || name != "Kyle" {
		t.Error("expected: Kyle, got:", name)
	}
}

func TestRegexMatcherWillReturnErrorForNoMatch(t *testing.T) {
	matcher := NewRegexMatcher("^(?i)(hi|hello|sup|hey), I'm (?P<name>(.*))$")
	msg := Message{Payload: "nonsense"}

	_, err := matcher.Match(msg)

	if err != ErrNoMatch {
		t.Error("expected: ErrNoMatch, got:", err)
	}
}

func TestRegexMatcherCanFindOneOutOfMany(t *testing.T) {
	matcher := NewRegexMatcher(
		"^(?i)hi, I'm (?P<name>(.*))$",
		"^testing (?P<thing>(.*))$")
	msg := Message{Payload: "testing something"}

	match, err := matcher.Match(msg)
	if err != nil {
		t.Error("expected: nil, got: ", err)
	}

	if name, ok := match["thing"]; !ok || name != "something" {
		t.Error("expected: Kyle, got:", name)
	}
}
