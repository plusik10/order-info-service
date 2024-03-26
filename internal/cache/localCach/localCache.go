package localcach

import (
	"fmt"
	"sync"
	"time"

	"github.com/plusik10/cmd/order-info-service/internal/cache"
)

var _ cache.Cache = (*localCache)(nil)

type localCache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]item
}

func NewLocalCache(defaultExpiration time.Duration, cleanupInterval time.Duration) cache.Cache {
	items := make(map[string]item)

	c := &localCache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	if cleanupInterval > 0 {
		go c.GC()
	}
	return c
}

// Count implements cache.Cache.
func (l *localCache) Count() int {
	return len(l.items)
}

// Delete implements cache.Cache.
func (l *localCache) Delete(key string) error {
	l.Lock()
	defer l.Unlock()

	if _, found := l.items[key]; !found {
		return fmt.Errorf("key %s not found", key)
	}
	delete(l.items, key)
	return nil
}

// Get implements cache.Cache.
func (l *localCache) Get(key string) (interface{}, bool) {
	l.RLock()
	defer l.RUnlock()
	item, found := l.items[key]

	if !found {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return item.Value, true
}

// Set implements cache.Cache.
func (l *localCache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64
	if duration == 0 {
		duration = l.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	l.Lock()
	defer l.Unlock()

	l.items[key] = item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

}

func (l *localCache) GC() {
	for {
		<-time.After(l.cleanupInterval)

		if l.items == nil {
			return
		}

		if keys := l.expiredKeys(); len(keys) != 0 {
			l.clearItems(keys)
		}
	}
}

func (l *localCache) expiredKeys() (keys []string) {
	l.RLock()

	defer l.RUnlock()

	for k, i := range l.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (l *localCache) clearItems(keys []string) {
	l.Lock()
	defer l.Unlock()

	for _, k := range keys {
		delete(l.items, k)
	}
}

type item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}
