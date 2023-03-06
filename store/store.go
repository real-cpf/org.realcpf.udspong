package store

import (
	"log"
	"sync"
)

type StoreAct struct {
	data sync.Map
}

func (s *StoreAct) Put(key []byte, value []byte) bool {
	k := string(key)
	log.Printf("store put %s \n", k)
	log.Printf("%p \n", &s.data)
	s.data.Store(k, string(value))
	return true
}

func (s *StoreAct) Get(key string) string {
	data, b := s.data.Load(key)
	log.Printf("%p \n", &s.data)
	log.Printf("store get %s %v \n", key, b)
	return data.(string)
}

func New() StoreAct {
	return StoreAct{
		data: sync.Map{},
	}
}
