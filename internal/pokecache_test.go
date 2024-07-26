package internal

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("what a path"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("couldn't find key %v", c.key)
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("the value from cache did not what was expected\nexpected: %v\nactual: %v", string(c.val), string(val))
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	const url = "https://example.com"
	cache := NewCache(baseTime)
	cache.Add(url, []byte("testdata"))

	_, ok := cache.Get(url)
	if !ok {
		t.Errorf("Expected to have cache entry but didn't find one")
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(url)
	if ok {
		t.Errorf("Expected to not find a cache entry but found one")
	}
}
