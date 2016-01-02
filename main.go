package main

import (
	"fmt"
	"github.com/pik4ez/drysh/plugin"
	_ "github.com/pik4ez/dryshplugin/outfile"
)

func main() {
	out := plugin.GetOut("outfile")
	result := out.Write("somestring")
	fmt.Println(result)
}
