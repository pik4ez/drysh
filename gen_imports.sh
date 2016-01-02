#!/bin/sh

if [ ! -f "$DRYSH_CFG" ]; then
    echo "ERROR: config must be specified via DRYSH_CFG env variable."
    echo "Usage: DRYSH_CFG=/path/to/config.toml go generate"
    exit 1;
fi

tpl_file=pluginimport/pluginimport.go.template
result_file=pluginimport/pluginimport.go

rm "$result_file"
cp "$tpl_file" "$result_file"

is_plugins_section=0
plugin_includes=''
while IFS='' read -r line || [[ -n "$line" ]]; do
    if [[ "$line" == "[plugins]" ]]; then
        is_plugins_section=1
    elif [[ "$line" == \[*\] ]]; then
        is_plugins_section=0
    fi
    if [[ $is_plugins_section -eq 1 ]]; then
        plugin_path=$(echo $line | cut -s -d'"' -f 2)
        if [[ "$plugin_path" != "" ]]; then
            plugin_includes="$plugin_includes\n_ \"$plugin_path\""
        fi
    fi
done < "$DRYSH_CFG"

plugin_includes=$(echo $plugin_includes | sed 's/\//\\\//g')
sed -i "s/\/\/ #PLUGIN_IMPORTS#/$plugin_includes/" "$result_file"
