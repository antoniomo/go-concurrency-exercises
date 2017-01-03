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
	d := newDinner()
	done := make(chan struct{})
	fmt.Println("Bon appétit!")
	wg.Add(4)
	go person("Alice", d, done, &wg)
	go person("Bob", d, done, &wg)
	go person("Carol", d, done, &wg)
	go person("Daniel", d, done, &wg)
	wg.Wait()
	fmt.Println("That was delicious!")
}

type dish struct {
	name string
	sync.Mutex
	morsels int
}

type dinner struct {
	dishes [5]*dish
	sync.Mutex
	total int
}

func newDinner() *dinner {
	d := [5]*dish{
		&dish{name: "chorizos", morsels: 5 + rand.Intn(6)},
		&dish{name: "chopitos", morsels: 5 + rand.Intn(6)},
		&dish{name: "croquetas", morsels: 5 + rand.Intn(6)},
		&dish{name: "patatas bravas", morsels: 5 + rand.Intn(6)},
		&dish{name: "pimientos de padrón", morsels: 5 + rand.Intn(6)},
	}
	total := 0
	for _, v := range d {
		total += v.morsels
	}
	fmt.Printf("Total morsels: %d\n", total)
	return &dinner{total: total, dishes: d}
}

func (d *dinner) getMorsel(done chan struct{}) string {
	var dd *dish
	for { // Inefficient as crap :D
		s := rand.Intn(5)
		dd = d.dishes[s]
		dd.Lock()
		if dd.morsels >= 1 {
			dd.morsels--
			dd.Unlock()
			break
		}
		dd.Unlock()
	}
	d.Lock()
	d.total--
	if d.total == 0 {
		close(done)
	}
	d.Unlock()
	return dd.name
}

func person(name string, d *dinner, done chan struct{}, wg *sync.WaitGroup) {
L:
	for {
		dishName := d.getMorsel(done)
		fmt.Printf("%s is enjoying some %s\n", name, dishName)
		// sleepTime := time.Duration(30+rand.Intn(151)) * time.Second // 30s to 3 min
		sleepTime := time.Duration(0+rand.Intn(1)) * time.Second
		time.Sleep(sleepTime)
		select {
		case <-done:
			break L
		default:
			continue
		}
	}
	wg.Done()
}
