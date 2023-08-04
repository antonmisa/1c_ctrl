package cache

import "time"

type item struct {
	Object     any
	Expiration *time.Time
}

func (i *item) Expired() bool {
	if i.Expiration == nil {
		return false
	}

	return i.Expiration.Before(time.Now())
}
