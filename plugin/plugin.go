package plugin

import (
	"github.com/pik4ez/drysh/out"
)

var (
	outs = make(map[string]out.Out)
)

func RegisterOut(name string, out out.Out) {
	if _, dup := outs[name]; dup {
		panic("plugin.out: Register called twice for out " + name)
	}
	outs[name] = out
}

func GetOut(name string) out.Out {
	return outs[name]
}
