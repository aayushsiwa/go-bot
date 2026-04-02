package services

import (
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

type cacheEntry struct {
	items     []*gofeed.Item
	expiresAt time.Time
}

var rssCache = struct {
	sync.RWMutex
	data map[string]cacheEntry
}{
	data: make(map[string]cacheEntry),
}
