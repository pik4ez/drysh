package main

import (
	"reflect"
	"testing"
	"time"
)

func TestOutput(t *testing.T) {
	outWaitTimeout := 50 * time.Millisecond

	inChan := make(chan Event)

	outChanOne := make(chan Event)
	outChanTwo := make(chan Event)
	outChanThree := make(chan Event)
	outChanFour := make(chan Event)

	router := Router{
		inChan:   inChan,
		outChans: make(map[string]chan Event),
		filters:  make(map[string]Filter),
	}

	router.RegisterOutput("out_one", outChanOne)
	router.RegisterOutput("out_two", outChanTwo)
	router.RegisterOutput("out_three", outChanThree)
	router.RegisterOutput("out_four", outChanFour)

	router.AddRuleOutput("tag_one", "out_one")
	router.AddRuleOutput("tag_two", "out_two")
	router.AddRuleOutput("**", "out_three")

	go router.Route()

	m := make(Msg)
	m["key_one"] = "value_one"
	ev := Event{Tag: "tag_one", Time: 12345, Msg: m}
	inChan <- ev
	expectOutput(t, outChanOne, ev, outWaitTimeout)

	m = make(Msg)
	m["key_two"] = "value_two"
	ev = Event{Tag: "tag_two", Time: 54321, Msg: m}
	inChan <- ev
	expectOutput(t, outChanTwo, ev, outWaitTimeout)

	m = make(Msg)
	m["key_three"] = "value_three"
	ev = Event{Tag: "tag_three", Time: 23432, Msg: m}
	inChan <- ev
	expectOutput(t, outChanThree, ev, outWaitTimeout)

	expectNoOutput(t, outChanFour, outWaitTimeout)
}

func TestFilter(t *testing.T) {
	outWaitTimeout := 50 * time.Millisecond

	inChan := make(chan Event)

	outChanOne := make(chan Event)

	filterOne := func(ev Event) Msg {
		ev.Msg["filtered_by"] = "filter_one"
		return ev.Msg
	}
	filterTwo := func(ev Event) Msg {
		ev.Msg["filtered_by"] = "filter_two"
		return ev.Msg
	}

	router := Router{
		inChan:   inChan,
		outChans: make(map[string]chan Event),
		filters:  make(map[string]Filter),
	}

	router.RegisterFilter("filter_one", filterOne)
	router.RegisterFilter("filter_two", filterTwo)
	router.RegisterOutput("out_one", outChanOne)

	router.AddRuleFilter("tag_one", "filter_one")
	router.AddRuleFilter("tag_two", "filter_two")
	router.AddRuleOutput("**", "out_one")

	go router.Route()

	m := make(Msg)
	m["key_one"] = "value_one"
	evSource := Event{Tag: "tag_one", Time: 12345, Msg: m}
	evExpected := evSource
	evExpected.Msg["filtered_by"] = "filter_one"
	inChan <- evSource
	expectOutput(t, outChanOne, evExpected, outWaitTimeout)

	m = make(Msg)
	m["key_two"] = "value_two"
	evSource = Event{Tag: "tag_two", Time: 54321, Msg: m}
	evExpected = evSource
	evExpected.Msg["filtered_by"] = "filter_two"
	inChan <- evSource
	expectOutput(t, outChanOne, evExpected, outWaitTimeout)

	m = make(Msg)
	m["key_three"] = "value_three"
	evSource = Event{Tag: "tag_three", Time: 23432, Msg: m}
	inChan <- evSource
	expectOutput(t, outChanOne, evSource, outWaitTimeout)
}

func expectOutput(t *testing.T, outChan chan Event, expectedEv Event, timeout time.Duration) {
	timeoutChan := make(chan bool, 1)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- true
	}()
	select {
	case <-timeoutChan:
		t.Fatalf("Event timed out, output: %v, expected event: %v\n", outChan, expectedEv)
	case resultEv := <-outChan:
		if reflect.DeepEqual(expectedEv, resultEv) {
			return
		}
		t.Fatalf("Events dont't match, output: %v, expected: %v, result: %v\n", outChan, expectedEv, resultEv)
	}
}

func expectNoOutput(t *testing.T, outChan chan Event, timeout time.Duration) {
	timeoutChan := make(chan bool, 1)
	go func() {
		time.Sleep(timeout)
		timeoutChan <- true
	}()
	select {
	case <-timeoutChan:
		return
	case resultEv := <-outChan:
		t.Fatalf("Unexpected event in output %v, event: %v\n", outChan, resultEv)
	}
}
