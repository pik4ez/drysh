package main

//go:generate ./gen_imports.sh

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/pik4ez/drysh/plugin"
	_ "github.com/pik4ez/drysh/pluginimport"
)

type out struct {
	Pattern string
	Plugin  string
	Config  toml.Primitive
}

type config struct {
	Outs []out `toml:"out"`
}

var configFile = "/home/vagrant/tmp/drysh_build/config.toml"

func main() {
	var cfg config
	var md toml.MetaData
	md, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, o := range cfg.Outs {
		out_plugin := plugin.GetOut(o.Plugin)
		out_cfg_prim := o.Config
		var plugin_conf = out_plugin.GetConfig()
		md.PrimitiveDecode(out_cfg_prim, plugin_conf)
	}

	out := plugin.GetOut("outfile")
	result := out.Write("somestring")
	fmt.Println(result)
}
