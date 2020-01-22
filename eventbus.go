package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type DataEvent struct {
	Data  interface{}
	Topic string
}

type DataChannel chan DataEvent

type DataChannels []DataChannel

type EventBus struct {
	subscribers map[string]DataChannels
	mux         sync.RWMutex
}

var eb = &EventBus{
	subscribers: map[string]DataChannels{},
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.mux.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.mux.Unlock()
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.mux.RLock()
	if chans, found := eb.subscribers[topic]; found {
		channels := append(DataChannels{}, chans...)
		go func(data DataEvent, dataChannels DataChannels) {
			for _, ch := range dataChannels {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}

	eb.mux.RUnlock()
}

func main() {
	ch1, ch2, ch3 := createChannel(), createChannel(), createChannel()
	eb.Subscribe("topic1", ch1)
	eb.Subscribe("topic2", ch2)
	eb.Subscribe("topic3", ch3)

	go publishTo("topic1", "Hi Topic 1")
	go publishTo("topic1", "Welcome to Topic 1")
	go publishTo("topic2", "Welcome to Topic 2")

	for {
		select {
		case d := <-ch1:
			go printDataEvent("ch1", d)
		case d := <-ch2:
			go printDataEvent("ch2", d)
		case d := <-ch3:
			go printDataEvent("ch3", d)
		}
	}

}

func createChannel() chan DataEvent {
	return make(chan DataEvent)
}

func printDataEvent(ch string, data DataEvent) {
	fmt.Printf("Channel: %s, Topic: %s; DataEvent: %v\n", ch, data.Topic, data.Data)
}

func publishTo(topic string, data string) {
	for {
		eb.Publish(topic, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}
