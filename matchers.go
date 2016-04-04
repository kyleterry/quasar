package quasar

import (
	"regexp"
)

type RegexMatcher struct {
	expressions []*regexp.Regexp
}

func NewRegexMatcher(expressions ...string) RegexMatcher {
	r := RegexMatcher{}
	for _, e := range expressions {
		r.expressions = append(r.expressions, regexp.MustCompile(e))
	}
	return r
}

func (r RegexMatcher) Match(msg Message) (Result, error) {
	res := make(Result)
	for _, expression := range r.expressions {
		match := expression.FindStringSubmatch(msg.Payload)
		if len(match) < 1 {
			continue
		}
		for i, name := range expression.SubexpNames() {
			res[name] = match[i]
		}
		return res, nil
	}
	return nil, ErrNoMatch
}
