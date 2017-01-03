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
	done := make(chan struct{})
	wg.Add(2)
	fmt.Println("Let's go for a walk!")
	go person("Alice", done, &wg)
	go person("Bob", done, &wg)
	wg.Wait()
	wg.Add(3)
	go alarm(&wg)
	wg.Wait()
}

func person(name string, done chan struct{}, wg *sync.WaitGroup) {
	fmt.Printf("%s started getting ready\n", name)
	sleepTime := time.Duration(60+rand.Intn(31)) * time.Second
	time.Sleep(sleepTime)
	wg.Done()
	fmt.Printf("%s spent %.0f seconds getting ready\n", name, sleepTime.Seconds())
	fmt.Printf("%s started putting on shoes\n", name)
	sleepTime = time.Duration(35+rand.Intn(11)) * time.Second
	time.Sleep(sleepTime)
	fmt.Printf("%s spent %.0f seconds putting on shoes\n", name, sleepTime.Seconds())
	select {
	case <-done:
		fmt.Println("Exiting and locking the door")
	default:
		close(done)
	}
	wg.Done()
}

func alarm(wg *sync.WaitGroup) {
	fmt.Println("Arming alarm")
	fmt.Println("Alarm is counting down")
	sleepTime := time.Duration(60) * time.Second
	time.Sleep(sleepTime)
	fmt.Println("Alarm is armed")
	wg.Done()
}
