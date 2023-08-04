package cache

import "time"

type clearJob struct {
	Interval time.Duration
	stop     chan struct{``}
}

func (j *clearJob) Run(c *cache) {
	j.stop = make(chan struct{})

	tick := time.Tick(j.Interval)

	for {
		select {
		case <-tick:
			c.DeleteExpired()
		case <-j.stop:
			return
		}
	}
}

func stopClearer(c *Cache) {
	c.clearer.stop <- struct{}{}
}

func runClearer(c *cache, ci time.Duration) {
	j := &clearJob{
		Interval: ci,
	}

	c.clearer = j

	go j.Run(c)
}
