package main

import (
	"sync"
)

func NewCache() *Cache{
	return &Cache{make(map[string]*MFile),sync.Mutex{}}
}

type MFile struct{
	str_image *string
	mx sync.Mutex
}
type Cache struct{
	data map[string] *MFile
	mx sync.Mutex
}

func (c *Cache)Data() (map[string] *MFile){
	return c.data
}

func (c *Cache) GetMFile(name string)*MFile{
	c.mx.Lock()
	defer c.mx.Unlock()
	if c.data[name] != nil{
		return c.data[name]
	}
	c.data[name] = &MFile{nil,sync.Mutex{}}
	return c.data[name]
}