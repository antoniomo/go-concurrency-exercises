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
	name    string
	morsels chan struct{}
}

type dinner [5]*dish

func newDinner() *dinner {
	d := dinner{
		&dish{name: "chorizos", morsels: make(chan struct{}, 5+rand.Intn(6))},
		&dish{name: "chopitos", morsels: make(chan struct{}, 5+rand.Intn(6))},
		&dish{name: "croquetas", morsels: make(chan struct{}, 5+rand.Intn(6))},
		&dish{name: "patatas bravas", morsels: make(chan struct{}, 5+rand.Intn(6))},
		&dish{name: "pimientos de padrón", morsels: make(chan struct{}, 5+rand.Intn(6))},
	}
	var empty struct{}
	for _, v := range d {
		for i := 0; i < cap(v.morsels); i++ {
			v.morsels <- empty
		}
	}
	return &d
}

func person(name string, d *dinner, done chan struct{}, wg *sync.WaitGroup) {
	var dishName string
L:
	for {
		select {
		case <-d[0].morsels:
			dishName = d[0].name
		case <-d[1].morsels:
			dishName = d[1].name
		case <-d[2].morsels:
			dishName = d[2].name
		case <-d[3].morsels:
			dishName = d[3].name
		case <-d[4].morsels:
			dishName = d[4].name
		default: // No more food :/
			break L
		}
		fmt.Printf("%s is enjoying some %s\n", name, dishName)
		// sleepTime := time.Duration(30+rand.Intn(151)) * time.Second // 30s to 3 min
		sleepTime := time.Duration(0+rand.Intn(1)) * time.Second
		time.Sleep(sleepTime)
	}
	wg.Done()
}
