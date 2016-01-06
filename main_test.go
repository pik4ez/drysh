package main

import (
	"reflect"
	"testing"
)

func TestMatcher(t *testing.T) {
	var testCases = []struct {
		pattern  string
		tag      string
		expected bool
	}{
		// Match all
		{"**", "", true},
		{"**", "some", true},
		{"**", "with.dot", true},
		// Partial match all
		{"a.**", "", false},
		{"a.**", "a", true},
		{"a.**", "a.b", true},
		{"a.**", "a.b.c", true},
		{"a.**", "b.a", false},
		{"a.**", "b.a.c", false},
		// Match single tag part
		{"*", "", true},
		{"*", "some", true},
		{"*", "with.dot", false},
		{"a.*", "", false},
		{"a.*", "a.b", true},
		{"a.*", "a.b.c", false},
		{"a.*", "b.a", false},
		{"a.*", "b.a.c", false},
		// Exact match
		{"the_tag", "", false},
		{"the_tag", "the_tag", true},
		{"the_tag", "other_tag", false},
		{"the_tag", "the_tag.suffix", false},
		{"the_tag", "prefix.the_tag", false},
		// Or pattern match
		{"{a,b}", "", false},
		{"{a,b}", "a", true},
		{"{a,b}", "b", true},
		{"{a,b}", "c", false},
		{"a b", "", false},
		{"a b", "a", true},
		{"a b", "b", true},
		{"a b", "c", false},

		// TODO
		// "**.a"
		// "a.**.b"
		// "*.a"
		// "a.*.b"
		// "a.{b,c}"
		// "a.{b,c}.*"
		// "a.{b,c}.**"
		// special chars
	}

	for _, tc := range testCases {
		m := Matcher{tc.pattern}
		expected := tc.expected
		result := m.Match(tc.tag)
		if expected != result {
			t.Errorf("Match tag \"%s\" against pattern \"%s\" failed, expected %v, got %v\n", tc.tag, tc.pattern, tc.expected, result)
		}
	}
}

func TestPatternExtractor(t *testing.T) {
	var testCases = []struct {
		raw      string
		expected []string
	}{
		{"", []string{}},
		{"*", []string{"*"}},
		{"**", []string{"**"}},
		{"a", []string{"a"}},
		{"a.*", []string{"a.*"}},
		{"a.**", []string{"a.**"}},
		{"a b", []string{"a", "b"}},
		{"a *", []string{"a", "*"}},
		{"a **", []string{"a", "**"}},
		{"{a}", []string{"a"}},
		{"{*}", []string{"*"}},
		{"{**}", []string{"**"}},
		{"{a,b}", []string{"a", "b"}},
		{"{a,*}", []string{"a", "*"}},
		{"{a,**}", []string{"a", "**"}},
		{"{a,}", []string{"a"}},
		{"a.{b,c}.**", []string{"a.b.**", "a.c.**"}},
	}

	for _, tc := range testCases {
		m := Matcher{tc.raw}
		expected := tc.expected
		result := m.ExtractPatterns()
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("Extractor failed on pattern \"%s\", expected %v, result %v", tc.raw, tc.expected, result)
		}
	}
}
