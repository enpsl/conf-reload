package base

import (
	"testing"
)

func TestLRUGetSet(t *testing.T) {
	l := CacheConstructor(2)
	tests := []struct {
		key string
		val interface{}
	}{
		{
			key: "hello",
			val: "world",
		},
		{
			key: "hello1",
			val: "world1",
		},
		{
			key: "hello1",
			val: "world1",
		},
		{
			key: "hello3",
			val: "world3",
		},
		{
			key: "hello4",
			val: "world4",
		},
		{
			key: "hello4",
			val: "world4",
		},
		{
			key: "hello4",
			val: "world4",
		},
	}
	for _, tc := range tests {
		l.Put(tc.key, tc.val)
		get, ok := l.Get(tc.key)
		if !ok || get != tc.val {
			t.Errorf("got=%s, want=%s", get, tc.val)
		}
	}
}

func TestLRUGetFlush(t *testing.T) {
	l := CacheConstructor(2)
	tests := []struct {
		key string
		val interface{}
	}{
		{
			key: "hello",
			val: "world",
		},
	}
	for _, tc := range tests {
		l.Flush()
		get, ok := l.Get(tc.key)
		if ok {
			t.Errorf("got=%s, want=%v", get, nil)
		}
	}
}
