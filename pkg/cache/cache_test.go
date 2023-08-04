package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Cache_Delete(t *testing.T) {
	type F func() *Cache

	tests := []struct {
		name string
		cf   F
		key  string
		want bool
	}{
		{
			name: "OK",
			cf: func() *Cache {
				c, _ := New(60 * time.Second)
				c.Set("1", 1)
				c.Set("2", 1)
				c.Set("3", 1)

				return c
			},
			key:  "1",
			want: true,
		},
		{
			name: "Error not found",
			cf: func() *Cache {
				c, _ := New(60 * time.Second)
				c.Set("2", 1)
				c.Set("3", 1)

				return c
			},
			key:  "1",
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := tt.cf()
			got := c.Delete(tt.key)

			require.Equal(t, got, tt.want)
		})
	}
}

func Test_Cache_DeleteExpired(t *testing.T) {
	type F func() *Cache

	tests := []struct {
		name string
		cf   F
		want int
	}{
		{
			name: "OK, long ttl",
			cf: func() *Cache {
				c, _ := New(60 * time.Second)
				c.Set("1", 1)
				c.Set("2", 1)
				c.Set("3", 1)

				return c
			},
			want: 3,
		},
		{
			name: "Ok, short ttl",
			cf: func() *Cache {
				c, _ := New(1 * time.Nanosecond)
				c.Set("1", 1)
				c.Set("2", 1)
				c.Set("3", 1)

				time.Sleep(2 * time.Nanosecond)

				return c
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := tt.cf()
			c.DeleteExpired()

			require.Equal(t, len(c.items), tt.want)
		})
	}
}

func Test_Cache_Flush(t *testing.T) {
	type F func() *Cache

	tests := []struct {
		name string
		cf   F
		want int
	}{
		{
			name: "OK, long ttl",
			cf: func() *Cache {
				c, _ := New(60 * time.Second)
				c.Set("1", 1)
				c.Set("2", 1)
				c.Set("3", 1)

				return c
			},
			want: 0,
		},
		{
			name: "Ok, short ttl",
			cf: func() *Cache {
				c, _ := New(1 * time.Nanosecond)
				c.Set("1", 1)
				c.Set("2", 1)
				c.Set("3", 1)

				return c
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := tt.cf()
			c.Flush()

			require.Equal(t, len(c.items), tt.want)
		})
	}
}

func Test_Cache_Get(t *testing.T) {
	type F func() *Cache

	tests := []struct {
		name  string
		cf    F
		key   string
		found bool
		want  any
	}{
		{
			name: "OK, long ttl",
			cf: func() *Cache {
				c, _ := New(60 * time.Second)
				c.Set("1", 1)
				c.Set("2", 1)
				c.Set("3", 1)

				return c
			},
			key:   "1",
			want:  1,
			found: true,
		},
		{
			name: "OK, modified",
			cf: func() *Cache {
				c, _ := New(60 * time.Second)
				c.Set("1", 1)
				c.Set("2", 1)
				c.Set("3", 1)
				c.Set("1", 2)

				return c
			},
			key:   "1",
			want:  2,
			found: true,
		},
		{
			name: "Not, short ttl expired",
			cf: func() *Cache {
				c, _ := New(1 * time.Nanosecond)
				c.Set("1", 1)
				c.Set("2", 1)
				c.Set("3", 1)

				time.Sleep(2 * time.Nanosecond)

				return c
			},
			key:   "1",
			want:  nil,
			found: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := tt.cf()
			v, found := c.Get(tt.key)

			require.Equal(t, found, tt.found)
			require.Equal(t, v, tt.want)
		})
	}
}

func Test_Cache_Set(t *testing.T) {
	type F func() *Cache

	tests := []struct {
		name string
		cf   F
		key  string
		want any
	}{
		{
			name: "OK",
			cf: func() *Cache {
				c, _ := New(60 * time.Second)

				return c
			},
			key:  "1",
			want: 1,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := tt.cf()
			c.Set(tt.key, tt.want)
			v, found := c.Get(tt.key)

			require.Equal(t, found, true)
			require.Equal(t, v, tt.want)
		})
	}
}
