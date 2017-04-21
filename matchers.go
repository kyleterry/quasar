package quasar

import (
	"regexp"
)

// RegexMatcher satisfies the Matcher interface and will match messages against
// a provided list of regular expressions.
type RegexMatcher struct {
	expressions []*regexp.Regexp
}

// NewRegexMatcher will return a new RegexMatcher setup with expressions
func NewRegexMatcher(expressions ...string) RegexMatcher {
	r := RegexMatcher{}
	for _, e := range expressions {
		r.expressions = append(r.expressions, regexp.MustCompile(e))
	}
	return r
}

func (r RegexMatcher) Match(msg Message) Result {
	res := make(Result)
	for _, expression := range r.expressions {
		match := expression.FindStringSubmatch(msg.Payload)
		if len(match) < 1 {
			continue
		}
		for i, name := range expression.SubexpNames() {
			res[name] = match[i]
		}
		return res
	}
	return nil
}
