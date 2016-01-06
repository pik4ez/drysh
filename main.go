package main

import (
	"fmt"
	"strings"
)

//go:generate ./gen_imports.sh

type Event struct {
	Tag  string
	Time int64
	Data map[string]string
}

type Matcher struct {
	RawPattern string
}

// Match checks if tag satisfies the pattern.
// TODO finish implementation (see matcher test)
func (m Matcher) Match(tag string) bool {
	matchAny := func(patterns []string, tag string) bool {
		matchOne := func(pattern string, tag string) bool {
			// Match all matches.
			if m.RawPattern == "**" {
				return true
			}
			// Exact match.
			if tag == m.RawPattern {
				return true
			}
			return false
		}

		for _, p := range patterns {
			if matchOne(p, tag) {
				return true
			}
		}
		return false
	}

	patterns := m.ExtractPatterns()
	return matchAny(patterns, tag)
}

// ExtractPatterns gets patterns from string.
// Supported syntax: whitespace separated ("a b c.*")
// and curly braces patterns limited to one closure
// ("{a,b,c.*}", "a.{b,c}.d).
func (m Matcher) ExtractPatterns() []string {
	// Patterns separated by spaces.
	if !strings.Contains(m.RawPattern, "{") {
		return strings.Fields(m.RawPattern)
	}

	// Pattern syntax {a,b}, a.{b,c}.* etc.
	var patterns []string
	idxOpen := strings.IndexRune(m.RawPattern, '{')
	idxClose := strings.IndexRune(m.RawPattern, '}')
	if idxClose == -1 {
		// Error: closing brace not found.
		return patterns
	}
	start := m.RawPattern[:idxOpen]
	middle := m.RawPattern[idxOpen+1 : idxClose]
	end := m.RawPattern[idxClose+1:]
	for _, m := range strings.Split(middle, ",") {
		if m == "" {
			continue
		}
		p := strings.Join([]string{start, m, end}, "")
		patterns = append(patterns, p)
	}
	return patterns
}

type Action struct {
	Type string
}

type Rule struct {
	Matcher Matcher
	Action  Action
}

func main() {
	var rules = []Rule{
		{Matcher{"tag_one"}, Action{"filter"}},
		{Matcher{"**"}, Action{"output"}},
	}
	fmt.Println(rules)
}
