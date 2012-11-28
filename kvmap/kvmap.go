package kvmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type KVMap struct {
	data     map[string]interface{}
	filename string
}

func New() *KVMap {
	return &KVMap{make(map[string]interface{}), ""}
}

func NewWith(key string, value interface{}) *KVMap {
	kv := New()
	kv.Set(key, value)
	return kv
}

func LoadJSONBuffer(data []byte) (*KVMap, error) {
	kv := New()
	if err := json.Unmarshal(data, &kv.data); err != nil {
		return nil, err
	}
	return kv, nil
}

func LoadJSONFile(name string) (*KVMap, error) {
	kv := New()
	kv.filename = name

	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	if err = dec.Decode(&kv.data); err != nil {
		return nil, err
	}

	return kv, nil
}

func (kv *KVMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(kv.data)
}

func (kv *KVMap) Merge(other *KVMap) {
	for key, value := range other.data {
		kv.data[key] = value
	}
	kv.filename = other.filename
}

func (kv *KVMap) errMsg(msg string) error {
	if kv.filename != "" {
		msg += fmt.Sprint(" in file ", kv.filename)
	}
	return errors.New("kvmap: " + msg)
}

func (kv *KVMap) errBadKey(key string) error {
	return kv.errMsg(fmt.Sprint("key ", key, " nonexistent"))
}

func (kv *KVMap) errBadType(key, typ string) error {
	return kv.errMsg(fmt.Sprint("non ", typ, " value with key ", key))
}

func (kv *KVMap) Set(key string, value interface{}) {
	kv.data[key] = value
}

func (kv *KVMap) get(key string) (interface{}, error) {
	value, ok := kv.data[key]
	if !ok {
		return nil, kv.errBadKey(key)
	}
	return value, nil
}

func (kv *KVMap) GetBool(key string) (bool, error) {
	v, err := kv.get(key)
	if err != nil {
		return false, err
	}

	val, ok := v.(bool)
	if !ok {
		return false, kv.errBadType(key, "boolean")
	}

	return val, nil
}

func (kv *KVMap) GetInt(key string) (int, error) {
	v, err := kv.get(key)
	if err != nil {
		return 0, err
	}

	val, ok := v.(int)
	if !ok {
		return 0, kv.errBadType(key, "integer")
	}

	return val, nil
}

func (kv *KVMap) GetString(key string) (string, error) {
	v, err := kv.get(key)
	if err != nil {
		return "", err
	}

	val, ok := v.(string)
	if !ok {
		return "", kv.errBadType(key, "string")
	}

	return val, nil
}
