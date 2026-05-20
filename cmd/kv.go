package gokv

import (
	"os"
	"strconv"
	"sync"
)

type KV struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// NewKV creates a new key-value store instance
func NewKV() *KV {
	return &KV{
		data: make(map[string][]byte),
	}
}

func (kv *KV) Set(key, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[string(key)] = val
	return nil
}
func (kv *KV) Delete(key []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.data, string(key))
	return nil
}

func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[string(key)]
	return val, ok
}

func (kv *KV) Auth(username, password string) bool {
	return username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD")
}
func (kv *KV) Exist(key []byte) (int, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	_, ok := kv.data[string(key)]
	if !ok {
		return 0, nil
	}
	return 1, nil
}

func (kv *KV) Incr(key []byte) (int, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	val, err := strconv.Atoi(string(kv.data[string(key)]))
	if err != nil {
		return 0, err
	}
	val++
	kv.data[string(key)] = []byte(strconv.Itoa(val))
	return val, nil
}

func (kv *KV) Decr(key []byte) (int, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	val, err := strconv.Atoi(string(kv.data[string(key)]))
	if err != nil {
		return 0, err
	}
	val--
	kv.data[string(key)] = []byte(strconv.Itoa(val))
	return val, nil
}
