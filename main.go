package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Customer struct {
	mutex sync.RWMutex
	id    string

	age int
}

func main() {
	wg := sync.WaitGroup{}
	var v uint64

	for i := 0; i < 3; i++ {
		go func() {
			wg.Add(1)
			atomic.AddUint64(&v, 1)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(v)
}

func (c *Customer) UpdateAge(age int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if age < 0 {
		return fmt.Errorf("age should be positive for customer %v", c)
	}

	c.age = age
	return nil
}

func (c *Customer) String() string {
	fmt.Println("enter string method")
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return fmt.Sprintf("id %s, age %d", c.id, c.age)
}
