package main

import (
	"time"
)

type Storage struct {
	data map[string]Value
}

type Value struct {
	value     string
	expiresAt time.Time
}

func (v Value) IsExpired() bool {
	if v.expiresAt.IsZero() {
		return false
	}

	return v.expiresAt.Before(time.Now())
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]Value),
	}
}

func (s *Storage) Set(key string, value string) {
	s.data[key] = Value{value: value}
}

func (s *Storage) SetWithExpiry(key string, value string, expiry time.Duration) {
	s.data[key] = Value{
		value:     value,
		expiresAt: time.Now().Add(expiry),
	}
}

func (s *Storage) Get(key string) (string, bool) {
	value, ok := s.data[key]

	if !ok {
		return "", false
	}

	if value.IsExpired() {
		delete(s.data, key)
		return "", false
	}

	return value.value, true
}
