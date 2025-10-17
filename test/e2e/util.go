package e2e

import (
	"regexp"
)

type RegexBuilder string

func (r RegexBuilder) String() string {
	return string(r)
}

func NewRegex() RegexBuilder {
	return ""
}

func (r RegexBuilder) Raw(regex string) RegexBuilder {
	return r + RegexBuilder(regex)
}

func (r RegexBuilder) Lit(lit string) RegexBuilder {
	return r.Raw(regexp.QuoteMeta(lit))
}

func (r RegexBuilder) Start() RegexBuilder {
	return r.Raw(`^`)
}

func (r RegexBuilder) End() RegexBuilder {
	return r.Raw(`$`)
}

func (r RegexBuilder) AnyTimes(regex RegexBuilder) RegexBuilder {
	return r.Raw(`(?:` + string(regex) + `)*`)
}

func (r RegexBuilder) OneOf(regexes ...RegexBuilder) RegexBuilder {
	r = r.Raw(`(?:`)
	for i, regex := range regexes {
		if i > 0 {
			r = r.Raw(`|`)
		}
		r = r.Raw(string(regex))
	}
	return r.Raw(`)`)
}

func (r RegexBuilder) OneOfLit(lits ...string) RegexBuilder {
	regexes := make([]RegexBuilder, 0, len(lits))
	for _, lit := range lits {
		regexes = append(regexes, NewRegex().Lit(lit))
	}
	return r.OneOf(regexes...)
}

func (r RegexBuilder) SeparatedByWhitespace(lits ...string) RegexBuilder {
	for i, lit := range lits {
		if i > 0 {
			r = r.Whitespace()
		}
		r = r.Lit(lit)
	}
	return r
}

func (r RegexBuilder) Datetime() RegexBuilder {
	// 2006-01-02 15:04:05 MST
	return r.Raw(`[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} [A-Z]+ `)
}

func (r RegexBuilder) FileSize() RegexBuilder {
	return r.Raw(`[0-9]+(?:\.[0-9]+)? (?:B|KB|MB|GB|TB)`)
}

func (r RegexBuilder) IPv4() RegexBuilder {
	// technically not correct but enough for our use case
	return r.Raw(`(?:[0-9]+\.){3}[0-9]+`)
}

func (r RegexBuilder) IPv6() RegexBuilder {
	return r.Raw(`[0-9a-f]+:[0-9a-f]+:[0-9a-f]+:[0-9a-f]+::[0-9]?`)
}

func (r RegexBuilder) AnyString() RegexBuilder {
	return r.Raw(`.+`)
}

func (r RegexBuilder) Identifier() RegexBuilder {
	return r.Raw(`[a-zA-Z0-9](?:[a-zA-Z0-9\-_.]*[a-zA-Z0-9])?`)
}

func (r RegexBuilder) Age() RegexBuilder {
	return r.Raw(`(?:just now|[0-9]+[smhd])`)
}

func (r RegexBuilder) HumanizeTime() RegexBuilder {
	return r.OneOf(`now`, `a long while ago`, `[0-9]+ (?:seconds?|minutes?|hours?|days?|months?|years?) (ago|from now)`)
}

func (r RegexBuilder) Whitespace() RegexBuilder {
	return r.Raw(`\s+`)
}

func (r RegexBuilder) OptionalWhitespace() RegexBuilder {
	return r.Raw(`\s*`)
}

func (r RegexBuilder) Newline() RegexBuilder {
	return r.Raw("\n")
}

func (r RegexBuilder) Int() RegexBuilder {
	return r.Raw(`[0-9]+`)
}

func (r RegexBuilder) Float() RegexBuilder {
	return r.Raw(`[0-9]*\.[0-9]*`)
}

func (r RegexBuilder) LocationName() RegexBuilder {
	return r.Raw(`[a-z]{3}[0-9]*`)
}

func (r RegexBuilder) Price() RegexBuilder {
	return r.Raw(`[€$] [0-9]+\.[0-9]+`)
}

func (r RegexBuilder) IBytes() RegexBuilder {
	return r.Raw(`[0-9]+(?:\.[0-9]+)? (?:B|[KMGTPE]iB)`)
}

func (r RegexBuilder) CountryCode() RegexBuilder {
	return r.Raw(`[A-Z]{2}`)
}
