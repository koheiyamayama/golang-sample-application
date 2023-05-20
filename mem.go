package main

import (
	"errors"

	"github.com/patrickmn/go-cache"
)

type MemDB struct {
	cache cache.Cache
}

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidDataType = errors.New("not allowed except string")
)

func NewMemDB() *MemDB {
	return &MemDB{
		cache: *cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (mdb *MemDB) List() map[string]string {
	items := mdb.cache.Items()
	retVal := make(map[string]string, 0)
	for key, item := range items {
		str, ok := item.Object.(string)
		if ok {
			retVal[key] = str
		}
	}

	return retVal
}

func (mdb *MemDB) Create(key string, value string) error {
	return mdb.cache.Add(key, value, cache.NoExpiration)
}

func (mdb *MemDB) Get(key string) (string, error) {
	val, found := mdb.cache.Get(key)
	if !found {
		return "", ErrNotFound
	}

	s, ok := val.(string)
	if ok {
		return s, nil
	} else {
		return "", ErrInvalidDataType
	}
}
