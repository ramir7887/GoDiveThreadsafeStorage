package kvstorage

import (
	"context"
	"errors"
	"log"
	"sync"
)

var (
	KeyNotExistsError  = errors.New("Key is not exists")
	ContextCancelError = errors.New("Cancel operation from context")
)

// KVStorage
type KVStorage interface {
	Get(context.Context, string) (interface{}, error)
	Put(context.Context, string, interface{}) error
	Delete(context.Context, string) error
}

// Storage потокобезопасное хранилище ключ-значение
type Storage struct {
	mu sync.RWMutex
	m  map[string]interface{}
}

// Get value if key exists else return error
func (s *Storage) Get(ctx context.Context, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		log.Println("Cancel from context")
		return nil, ContextCancelError
	default:
		log.Println("Getting value")
	}

	if _, ok := s.m[key]; !ok {
		return nil, KeyNotExistsError
	}
	return s.m[key], nil
}

// Put add key and value
func (s *Storage) Put(ctx context.Context, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		log.Println("Cancel from context")
		return ContextCancelError
	default:
		log.Println("Put value", value)
	}

	if s.m == nil {
		s.m = make(map[string]interface{})
	}
	s.m[key] = value
	return nil
}

// Delete
func (s *Storage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		log.Println("Cancel from context")
		return ContextCancelError
	default:
		log.Println("Delete value")
	}

	if _, ok := s.m[key]; !ok {
		return KeyNotExistsError
	}
	delete(s.m, key)
	return nil
}
