package cache

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		ttl time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *Cache
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Cache_Delete(t *testing.T) {
	type fields struct {
		ttl time.Duration
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New() & Cache{
				Mutex:   tt.fields.Mutex,
				ttl:     tt.fields.ttl,
				items:   tt.fields.items,
				clearer: tt.fields.clearer,
			}
			if got := c.Delete(tt.args.key); got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_DeleteExpired(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		ttl     time.Duration
		items   map[string]*item
		clearer *clearJob
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				Mutex:   tt.fields.Mutex,
				ttl:     tt.fields.ttl,
				items:   tt.fields.items,
				clearer: tt.fields.clearer,
			}
			c.DeleteExpired()
		})
	}
}

func Test_cache_Flush(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		ttl     time.Duration
		items   map[string]*item
		clearer *clearJob
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				Mutex:   tt.fields.Mutex,
				ttl:     tt.fields.ttl,
				items:   tt.fields.items,
				clearer: tt.fields.clearer,
			}
			c.Flush()
		})
	}
}

func Test_cache_Get(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		ttl     time.Duration
		items   map[string]*item
		clearer *clearJob
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				Mutex:   tt.fields.Mutex,
				ttl:     tt.fields.ttl,
				items:   tt.fields.items,
				clearer: tt.fields.clearer,
			}
			got, got1 := c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_cache_Set(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		ttl     time.Duration
		items   map[string]*item
		clearer *clearJob
	}
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				Mutex:   tt.fields.Mutex,
				ttl:     tt.fields.ttl,
				items:   tt.fields.items,
				clearer: tt.fields.clearer,
			}
			c.Set(tt.args.key, tt.args.value)
		})
	}
}

func Test_cache_delete(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		ttl     time.Duration
		items   map[string]*item
		clearer *clearJob
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				Mutex:   tt.fields.Mutex,
				ttl:     tt.fields.ttl,
				items:   tt.fields.items,
				clearer: tt.fields.clearer,
			}
			c.delete(tt.args.key)
		})
	}
}

func Test_cache_get(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		ttl     time.Duration
		items   map[string]*item
		clearer *clearJob
	}
	type args struct {
		k string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   any
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				Mutex:   tt.fields.Mutex,
				ttl:     tt.fields.ttl,
				items:   tt.fields.items,
				clearer: tt.fields.clearer,
			}
			got, got1 := c.get(tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_cache_set(t *testing.T) {
	type fields struct {
		Mutex   sync.Mutex
		ttl     time.Duration
		items   map[string]*item
		clearer *clearJob
	}
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				Mutex:   tt.fields.Mutex,
				ttl:     tt.fields.ttl,
				items:   tt.fields.items,
				clearer: tt.fields.clearer,
			}
			c.set(tt.args.key, tt.args.value)
		})
	}
}

func Test_newCache(t *testing.T) {
	type args struct {
		ttl time.Duration
	}
	tests := []struct {
		name string
		args args
		want *cache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCache(tt.args.ttl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCache() = %v, want %v", got, tt.want)
			}
		})
	}
}
