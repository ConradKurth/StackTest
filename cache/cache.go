package cache

import (
	"sync"
	"time"
)

type Item struct {
	Key   string
	Value interface{}

	Created time.Time
}

var c []Item
var cLock sync.Mutex

func AddItem(k string, v interface{}) {
	cLock.Lock()
	defer cLock.Unlock()
	for i := 0; i < len(c); i++ {
		if c[i].Key == k {
			c[i].Value = v
			c[i].Created = time.Now().UTC()
			return
		}
	}
	c = append(c, Item{Key: k, Value: v, Created: time.Now().UTC()})
}

func GetItem(k string) interface{} {
	cLock.Lock()
	cLock.Unlock()
	for i := len(c) - 1; i >= 0; i-- {
		if c[i].Key == k {
			d := time.Now().UTC().Sub(c[i].Created)
			if d < time.Minute*5 {
				return c[i].Value
			} else {
				c = append(c[:i], c[i+1:]...)
			}
		}
	}
	return nil
}
