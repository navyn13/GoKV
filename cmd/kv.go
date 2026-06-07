package gokv

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type KV struct {
	mu   sync.RWMutex
	data map[string]value
}
type value struct {
	data   string
	expire time.Time
}

// NewKV creates a new key-value store instance
func NewKV() *KV {
	return &KV{
		data: make(map[string]value),
	}
}

func (kv *KV) Set(key, content []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	value := value{
		data:   string(content),
		expire: time.Time{},
	}
	kv.data[string(key)] = value
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
	if !ok {
		return []byte{}, false
	}
	if val.expire.IsZero() {
		return []byte(val.data), true
	}
	if val.expire.Before(time.Now()) {
		delete(kv.data, string(key))
		return []byte{}, false
	}
	return []byte(val.data), true
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

	val, err := strconv.Atoi(kv.data[string(key)].data)
	if err != nil {
		return 0, err
	}
	val++
	kv.data[string(key)] = value{
		data:   strconv.Itoa(val),
		expire: time.Time{},
	}
	return val, nil
}

func (kv *KV) Decr(key []byte) (int, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	val, err := strconv.Atoi(kv.data[string(key)].data)
	if err != nil {
		return 0, err
	}
	val--
	kv.data[string(key)] = value{
		data:   strconv.Itoa(val),
		expire: time.Time{},
	}
	return val, nil
}

func (kv *KV) Expire(key []byte, seconds int) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	_, ok := kv.data[string(key)]
	if !ok {
		return fmt.Errorf("key not found")
	}
	kv.data[string(key)] = value{
		data:   kv.data[string(key)].data,
		expire: time.Now().Add(time.Duration(seconds) * time.Second),
	}
	return nil
}
