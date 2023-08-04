package cache

import "time"

type clearJob struct {
	interval time.Duration
	stop     chan struct{}
}

func (j *clearJob) Run(c *cache) {
	ticker := time.NewTicker(j.interval)

	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopClearer(c *Cache) {
	c.clearer.stop <- struct{}{}
}

func runClearer(c *cache, ci time.Duration) {
	j := &clearJob{
		interval: ci,
		stop:     make(chan struct{}),
	}

	c.clearer = j

	go j.Run(c)
}
