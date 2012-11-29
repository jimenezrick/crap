package main

import "crap/kvmap"

var configFiles []string = []string{"/etc/crap/conf.json", "conf.json"}

func defaultConfig() *kvmap.KVMap {
	config := kvmap.New()
	config.Set("log.debug", true)
	config.Set("log.syslog", true)
	config.Set("store.path", "/tmp")
	config.Set("store.dir_permissions", 0700)
	config.Set("store.file_permissions", 0600)
	config.Set("network.listen_address", ":9000")
	return config
}

func mergeConfigFile(config *kvmap.KVMap, name string) {
	configFile, err := kvmap.LoadJSONFile(name)
	if err == nil {
		config.Merge(configFile)
	}
}

func loadConfigFiles(config *kvmap.KVMap, names []string) {
	for _, name := range names {
		mergeConfigFile(config, name)
	}
}

func loadConfig() *kvmap.KVMap {
	config := defaultConfig()
	loadConfigFiles(config, configFiles)
	return config
}
