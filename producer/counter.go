package producer

import "sync"

type Counter struct {
	i int
	sync.Mutex
}

func (c *Counter) Increment() int {
	c.Lock()
	defer c.Unlock()
	c.i++
	return c.i
}

func (c *Counter) Read() int {
	c.Lock()
	defer c.Unlock()
	return c.i
}