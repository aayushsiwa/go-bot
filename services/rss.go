package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

//var FeedMap = map[string]map[string]string{
//	"tech": {
//		"tech": "https://hnrss.org/frontpage",
//	},
//	"gaming": {"feedburner": "https://feeds.feedburner.com/ign/all"},
//}

var FeedMap = map[string]string{
	"tech":   "https://hnrss.org/frontpage",
	"news":   "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml",
	"gaming": "https://feeds.feedburner.com/ign/all",
}

func FetchLatest(url string, limit int) ([]*gofeed.Item, error) {
	parser := gofeed.NewParser()

	feed, err := parser.ParseURL(url)
	if err != nil {
		return nil, err
	}

	if len(feed.Items) > limit {
		return feed.Items[:limit], nil
	}

	return feed.Items, nil
}

func GetCachedFeed(url string, ttl time.Duration) ([]*gofeed.Item, error) {
	now := time.Now()

	rssCache.RLock()
	entry, found := rssCache.data[url]
	rssCache.RUnlock()

	// ✅ return cached if valid
	if found && now.Before(entry.expiresAt) {
		return entry.items, nil
	}

	// ❌ fetch fresh
	items, err := FetchLatest(url, Cfg.RSSFeedSize)
	if err != nil {
		return nil, err
	}

	// ✅ store in cache
	rssCache.Lock()
	rssCache.data[url] = cacheEntry{
		items:     items,
		expiresAt: now.Add(ttl),
	}
	rssCache.Unlock()

	return items, nil
}

func GetFeedItems(feed string) ([]*gofeed.Item, error) {
	url, ok := FeedMap[feed]
	if !ok {
		return nil, errors.New("invalid feed. Options: tech, gaming")
	}

	items, err := GetCachedFeed(url, 5*time.Minute)
	if err != nil || len(items) == 0 {
		return nil, fmt.Errorf("failed to fetch")
	}

	return items, nil
}
