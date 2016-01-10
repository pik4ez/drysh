package main

//go:generate ./gen_imports.sh

type Msg map[string]string

type Event struct {
	Tag  string
	Time int64
	Msg  Msg
}

func main() {
	//
}
