package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Handler func(a, b int)

type Publisher interface {
	Publish(topic string, a, b int)
}

type Subscriber interface {
	Subscribe(topic string, fn Handler)
}

type EventBus struct {
	mutex       sync.Mutex
	subscribers map[string][]Handler
	Publisher
	Subscriber
}

func New() *EventBus {
	return &EventBus{
		mutex:       sync.Mutex{},
		subscribers: make(map[string][]Handler),
	}
}

func (bus *EventBus) Subscribe(topic string, fn Handler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	if _, found := bus.subscribers[topic]; !found {
		bus.subscribers[topic] = make([]Handler, 0)
	}

	bus.subscribers[topic] = append(bus.subscribers[topic], fn)
}

func (bus *EventBus) Publish(topic string, a, b int) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	if handlers, found := bus.subscribers[topic]; found {
		go func(a, b int) {
			for _, handler := range handlers {
				handler(a, b)
			}
		}(a, b)
	}
}

func add(a, b int) {
	log.Println("add", a+b)
}

func sub(a, b int) {
	log.Println("sub", a-b)
}

func main() {
	var (
		wg     sync.WaitGroup
		topics = map[int]string{
			0: "add",
			1: "sub",
		}
	)

	rand.Seed(time.Now().UTC().UnixNano())

	bus := New()
	bus.Subscribe("add", add)
	bus.Subscribe("sub", sub)

	wg.Add(1)

	go func() {
		for {
			bus.Publish(topics[rand.Intn(2)], rand.Intn(10), rand.Intn(10))
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()

	wg.Wait()
}
