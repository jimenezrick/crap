package config

import (
	"fmt"
	"os"
)

import "crap/kvmap"

func defaultConfig() *KVMap {
	config := New()

	config.Set("log.debug", true)
	config.Set("log.syslog", true)
	config.Set("store.path", ".")
	config.Set("store.dir_permissions", "0700")
	config.Set("store.file_permissions", "0600")
	config.Set("network.listen_address", ":9000")

	return config
}

func mergeConfigFile(config *KVMap, name string) {
	if kvmap, err := kvmap.LoadJSONFile(name); err == nil {
		config.Merge(kvmap)
		fmt.Println("Config file", name, "loaded") // TODO: Check verbose flag
	} else if !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Syntax error in config file", name+":", err)
		os.Exit(1)
	}
}

func LoadConfig(names []string) Config {
	config := defaultConfig()
	for _, name := range names {
		mergeConfigFile(config, name)
	}
	return config
}
