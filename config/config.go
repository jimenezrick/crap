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

type KVMap struct {
	kvmap.KVMap
}

func New() *KVMap {
	return &KVMap{*kvmap.New()}
}

func (kv KVMap) GetBool(key string) bool {
	b, err := kv.KVMap.GetBool(key)
	if err != nil {
		panic(ConfigError{err})
	}
	return b
}

func (kv KVMap) GetInt(key string) int {
	i, err := kv.KVMap.GetInt(key)
	if err != nil {
		panic(ConfigError{err})
	}
	return i
}

func (kv KVMap) GetString(key string) string {
	s, err := kv.KVMap.GetString(key)
	if err != nil {
		panic(ConfigError{err})
	}
	return s
}

func (kv KVMap) GetIntString(key string) int {
	i, err := kv.KVMap.GetIntString(key)
	if err != nil {
		panic(ConfigError{err})
	}
	return i
}
