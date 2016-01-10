package main

import (
	"fmt"
	"time"
)

type ActionType int

const (
	ACTION_TYPE_FILTER ActionType = iota
	ACTION_TYPE_OUTPUT

	IDLE_TIMEOUT = 50 * time.Millisecond
)

type Action struct {
	Type ActionType
	Name string
}

type Rule struct {
	Matcher Matcher
	Action  Action
}

type Filter func(Event) Msg

type Router struct {
	inChan   chan Event
	outChans map[string]chan Event
	filters  map[string]Filter
	rules    []Rule
}

func (r Router) Route() {
	for {
		select {
		case ev := <-r.inChan:
			actions := r.getActionsSequence(ev)
			for _, action := range actions {
				switch action.Type {
				case ACTION_TYPE_FILTER:
					filter := r.filters[action.Name]
					ev.Msg = filter(ev)
				case ACTION_TYPE_OUTPUT:
					r.outChans[action.Name] <- ev
				default:
					panic(fmt.Sprintf("Unknown type in action %v", action))
				}
			}
		default:
			time.Sleep(IDLE_TIMEOUT)
		}
	}
}

func (r Router) SetInChan(inChan chan Event) {
	r.inChan = inChan
}

func (r *Router) AddRuleFilter(tag string, name string) {
	rule := Rule{Matcher{tag}, Action{ACTION_TYPE_FILTER, name}}
	r.rules = append(r.rules, rule)
}

func (r *Router) AddRuleOutput(tag string, name string) {
	rule := Rule{Matcher{tag}, Action{ACTION_TYPE_OUTPUT, name}}
	r.rules = append(r.rules, rule)
}

func (r Router) RegisterFilter(name string, filter Filter) {
	r.filters[name] = filter
}

func (r Router) RegisterOutput(name string, outChan chan Event) {
	r.outChans[name] = outChan
}

func (r Router) getActionsSequence(ev Event) []Action {
	var actions []Action
	for _, r := range r.rules {
		if r.Matcher.Match(ev.Tag) {
			actions = append(actions, r.Action)
			if r.Action.Type == ACTION_TYPE_OUTPUT {
				return actions
			}
		}
	}
	return actions
}
