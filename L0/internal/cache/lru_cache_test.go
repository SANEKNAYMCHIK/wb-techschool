package cache

import (
	"testing"
)

func TestNewLRUCache(t *testing.T) {
	t.Run("zero capacity", func(t *testing.T) {
		cache := NewLRUCache(0)
		if cache == nil {
			t.Fatal("Cache should not be nil")
		}
		if cache.list.Len() != 0 {
			t.Errorf("Expected 0 items, got %d", cache.list.Len())
		}
	})

	t.Run("normal capacity", func(t *testing.T) {
		cache := NewLRUCache(10)
		if cache == nil {
			t.Fatal("Cache should not be nil")
		}
	})
}

func TestLRUCache_SetGet(t *testing.T) {
	cache := NewLRUCache(2)

	t.Run("set and get", func(t *testing.T) {
		cache.Set("key1", "value1")
		if val, ok := cache.Get("key1"); !ok || val != "value1" {
			t.Errorf("Get failed, expected value1, got %v", val)
		}
	})

	t.Run("get non-existent", func(t *testing.T) {
		if _, ok := cache.Get("missing"); ok {
			t.Error("Should not find missing key")
		}
	})
}

func TestLRUCache_Removing(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	t.Run("all keys present", func(t *testing.T) {
		if _, ok := cache.Get("key1"); !ok {
			t.Error("key1 should be present")
		}
		if _, ok := cache.Get("key2"); !ok {
			t.Error("key2 should be present")
		}
		if _, ok := cache.Get("key3"); !ok {
			t.Error("key3 should be present")
		}
	})

	t.Run("removing when adding new", func(t *testing.T) {
		cache.Set("key4", "value4")
		if _, ok := cache.Get("key1"); ok {
			t.Error("key1 should be removing")
		}
		if _, ok := cache.Get("key4"); !ok {
			t.Error("key4 should be present")
		}
	})

	t.Run("access affecting on removing", func(t *testing.T) {
		cache.Get("key2")

		cache.Set("key5", "value5")

		if _, ok := cache.Get("key3"); ok {
			t.Error("key3 should be removing")
		}
		if _, ok := cache.Get("key2"); !ok {
			t.Error("key2 should still be present")
		}
	})
}

func TestLRUCache_UpdateValue(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Set("key1", "value1")
	cache.Set("key1", "updated")

	if val, ok := cache.Get("key1"); !ok || val != "updated" {
		t.Errorf("Update failed, expected 'updated', got '%v'", val)
	}
}

func TestLRUCache_Len(t *testing.T) {
	cache := NewLRUCache(3)

	if cache.list.Len() != 0 {
		t.Errorf("Expected length 0, got %d", cache.list.Len())
	}

	cache.Set("key1", "value1")
	if cache.list.Len() != 1 {
		t.Errorf("Expected length 1, got %d", cache.list.Len())
	}

	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	if cache.list.Len() != 3 {
		t.Errorf("Expected length 3, got %d", cache.list.Len())
	}

	cache.Set("key4", "value4")
	if cache.list.Len() != 3 {
		t.Errorf("Expected length 3 after removing, got %d", cache.list.Len())
	}
}
