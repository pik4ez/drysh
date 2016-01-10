package main

import "fmt"

//go:generate ./gen_imports.sh

type Event struct {
	Tag  string
	Time int64
	Data map[string]string
}

type ActionType int8

const (
	ACTION_TYPE_FILTER ActionType = iota
	ACTION_TYPE_OUTPUT
)

type Action struct {
	Type ActionType
}

type Rule struct {
	Matcher Matcher
	Action  Action
}

func main() {
	var rules = []Rule{
		{Matcher{"tag_one"}, Action{ACTION_TYPE_FILTER}},
		{Matcher{"**"}, Action{ACTION_TYPE_OUTPUT}},
	}
	fmt.Println(rules)
}
