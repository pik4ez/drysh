#!/bin/sh

tpl_file=pluginimport/pluginimport.go.template
result_file=pluginimport/pluginimport.go
rm "$result_file"
cp "$tpl_file" "$result_file"
sed -i 's/\/\/ #PLUGIN_IMPORTS#/_ "github.com\/pik4ez\/dryshplugin\/outfile"/' "$result_file"
