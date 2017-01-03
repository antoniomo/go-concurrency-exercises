package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	c := newiCafe()
	wg.Add(25)
	for i := 0; i < 25; i++ {
		go tourist(i, c, &wg)
		time.Sleep(1 * time.Millisecond)
	}
	wg.Wait()
}

type iCafe struct {
	freeComputers chan bool
}

func newiCafe() *iCafe {
	c := &iCafe{make(chan bool, 8)}
	for i := 0; i < 8; i++ {
		c.freeComputers <- true
	}
	return c
}

func tourist(number int, c *iCafe, wg *sync.WaitGroup) {
	doOnceIfBlock := make(chan struct{})
	close(doOnceIfBlock)
L:
	for {
		// First non-blocking select simulates priority the first time
		select {
		case <-c.freeComputers:
			break L
		default:
		}
		select {
		case <-c.freeComputers:
			break L
		case <-doOnceIfBlock:
			fmt.Printf("Tourist %d waiting for turn.\n", number)
			doOnceIfBlock = nil
		}
	}
	fmt.Printf("Tourist %d is online.\n", number)
	sleepTime := time.Duration(5+rand.Intn(5)) * time.Second
	time.Sleep(sleepTime)
	c.freeComputers <- true
	fmt.Printf("Tourist %d is done.\n", number)
	wg.Done()
}
