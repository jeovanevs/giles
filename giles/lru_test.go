package main

import (
	"testing"
)

func TestGet(t *testing.T) {
	lru := NewLRU(uint32(4))

	val, ok := lru.Get("asdf")
	if ok != false {
		t.Error("ok should be false but is", ok)
	}
	lru.Set("asdf", "asdfvalue")
	val, ok = lru.Get("asdf")
	if val != "asdfvalue" {
		t.Error("LRU.Get does not return correct value")
	}

	if lru.cache["asdf"] != "asdfvalue" {
		t.Error("LRU.cache does not contain key/value after Get")
	}

	if lru.queue.Front().Value.(string) != "asdf" {
		t.Error("LRU.queue does not contain key: asdf")
	}
}

func TestEviction(t *testing.T) {
	lru := NewLRU(2)
	val1, ok := lru.Get("a")
	if ok != false {
		t.Error("ok should be false")
	}
	lru.Set("a", "avalue")
	val1, ok = lru.Get("a")
	if ok != true {
		t.Error("ok should be true")
	}
	if val1 != "avalue" {
		t.Error("lru.Get: val1 should be avalue but is", val1)
	}

	lru.Set("b", "bvalue")
	lru.Set("c", "cvalue")

	if len(lru.cache) != 2 {
		t.Error("lru.Cache size should be 2 but is", len(lru.cache))
	}
	if lru.queue.Len() != 2 {
		t.Error("lru.queue len should be 2 but is", lru.queue.Len())
	}
	if lru.queue.Front().Value.(string) != "c" {
		t.Error("Most recently used item should be 'c' but is", lru.queue.Front().Value.(string))
	}
	if lru.queue.Back().Value.(string) != "b" {
		t.Error("Most recently used item should be 'b' but is", lru.queue.Back().Value.(string))
	}
}