package main

import "fmt"

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Set(key string, value string) {
	s.data[key] = value
}

func (s *Storage) Get(key string) (string, error) {
	value, ok := s.data[key]

	if !ok {
		return "", fmt.Errorf("key does not exist: %v", key)
	}

	return value, nil
}
