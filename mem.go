package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/oklog/ulid/v2"
	"github.com/patrickmn/go-cache"
)

type MemDB struct {
	cache *cache.Cache
}

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidDataType = errors.New("not allowed except string")
)

func NewMemDB() *MemDB {
	b, err := os.ReadFile("./posts.json")
	if err != nil {
		return &MemDB{
			cache: cache.New(cache.NoExpiration, cache.NoExpiration),
		}
	} else {
		items := map[string]cache.Item{}
		v := []*Post{}
		_ = json.Unmarshal(b, &v)

		for _, post := range v {
			id := ulid.Make().String()
			post.ID = id
			items[id] = cache.Item{Object: post.String(), Expiration: 0}
		}

		return &MemDB{
			cache: cache.NewFrom(cache.NoExpiration, cache.NoExpiration, items),
		}
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
