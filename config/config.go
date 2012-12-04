package config

import "crap/kvmap"

type Config interface {
	GetBool(key string) bool
	GetInt(key string) int
	GetString(key string) string
	GetIntString(key string) int
}

type ConfigError struct {
	Err error
}

type kvMap struct {
	kvmap.KVMap
}

func New() *kvMap {
	return &kvMap{*kvmap.New()}
}

func (kv kvMap) GetBool(key string) bool {
	b, err := kv.KVMap.GetBool(key)
	if err != nil {
		panic(ConfigError{err})
	}
	return b
}

func (kv kvMap) GetInt(key string) int {
	i, err := kv.KVMap.GetInt(key)
	if err != nil {
		panic(ConfigError{err})
	}
	return i
}

func (kv kvMap) GetString(key string) string {
	s, err := kv.KVMap.GetString(key)
	if err != nil {
		panic(ConfigError{err})
	}
	return s
}

func (kv kvMap) GetIntString(key string) int {
	i, err := kv.KVMap.GetIntString(key)
	if err != nil {
		panic(ConfigError{err})
	}
	return i
}
