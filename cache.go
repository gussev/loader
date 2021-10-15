package main

import (
	"sync"
)

func NewCache() *Cache{
	return &Cache{make(map[string]*string),sync.Mutex{}}
}

type Cache struct{
	data map[string] *string
	mx sync.Mutex
}

func (c *Cache)Get(key string,convert func(name string) ( *string, error)) (*string,error){
	c.mx.Lock()
	defer c.mx.Unlock()
	str_image := c.data[key]
	if str_image != nil{
		return str_image,nil
	}
	str_image,err := convert(key)
	if err != nil{
		return nil,err
	}
	c.data[key] = str_image
	return str_image,nil
}