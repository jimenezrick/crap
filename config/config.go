package config

var conf map[string]interface{} = make(map[string]interface{})

func Set(key string, value interface{}) {
	conf[key] = value
}

func get(key string) interface{} {
	value, ok := conf[key]
	if !ok {
		panic("crap/config: non existent key")
	}
	return value
}

func GetInt(key string) int {
	value, ok := get(key).(int)
	if !ok {
		panic("crap/config: non integer value")
	}
	return value
}

func GetString(key string) string {
	value, ok := get(key).(string)
	if !ok {
		panic("crap/config: non string value")
	}
	return value
}
