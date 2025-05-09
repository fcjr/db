package store

import "sync"

type Store struct {
	kv     map[string]string
	kvLock *sync.RWMutex
}

func New() *Store {
	return &Store{
		kv:     map[string]string{},
		kvLock: &sync.RWMutex{},
	}
}

func (s *Store) Get(key string) string {
	s.kvLock.RLock()
	defer s.kvLock.RUnlock()

	return s.kv[key]
}

func (s *Store) Set(key, val string) {
	s.kvLock.Lock()
	defer s.kvLock.Unlock()

	s.kv[key] = val
}
