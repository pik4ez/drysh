package main

//go:generate ./gen_imports.sh

import (
	"fmt"
	"github.com/pik4ez/drysh/plugin"
	_ "github.com/pik4ez/drysh/pluginimport"
)

func main() {
	out := plugin.GetOut("outfile")
	result := out.Write("somestring")
	fmt.Println(result)
}
