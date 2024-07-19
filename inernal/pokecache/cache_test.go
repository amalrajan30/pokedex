package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T){
	const interval = 5 * time.Second

	testCases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
			val: []byte("testData"),
		},
		{
			key: "https://pokeapi.co/api/v2/location-area/?offset=40&limit=20",
			val: []byte("nextTestData"),
		},
	}

	for i, tc := range testCases {
		fmt.Printf("Test case %v \n", i)
		cache := NewCache(interval)

		cache.Add(tc.key, tc.val)

		val, ok := cache.Get(tc.key)

		if !ok {
			t.Fatalf("expected to get val of %s", tc.key)
		}

		if string(val) != string(tc.val) {
			t.Fatalf("expected the values to match: %v, %v", val, tc.val)
		}
	}

	cache := NewCache(interval)

	_, ok := cache.Get("not cached")
	if ok {
		t.Fatalf("expected key not found")
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 10 * time.Millisecond

	cache := NewCache(baseTime)
	cacheKey := "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20"

	cache.Add(cacheKey, []byte("testdata"))

	_, ok := cache.Get(cacheKey)
	if !ok {
		t.Fatalf("expected to find key")
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(cacheKey)
	if ok {
		t.Fatalf("expected the key to expire")
	}
}