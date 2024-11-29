package wh03

import (
	"sync"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type IPubSub interface {
	Publish(topic string, msg int) error
	Subscribe(topic string) <-chan int
	Close()
	Init(bufferSize int)
}

type PubSub struct {
	mu         sync.RWMutex
	submap     map[string][]chan int
	closed     bool
	bufferSize int
}

func (it *PubSub) Publish(topic string, msg int) error {
	logw.Debugf("^PubSub.Publish - topic: %s, msg: %08b", topic, msg)
	defer logw.Debugf("$PubSub.Publish")

	it.mu.RLock()
	defer it.mu.RUnlock()
	logw.Debug("PubSub.Publish - acquired mutex")
	defer logw.Debug("PubSub.Publish - releasing mutex")

	logw.Infof("PubSub.Publish - topic: %s, msg: %08b", topic, msg)

	t := it.submap[topic]

	for _, sub := range t {
		sub <- msg
	}

	return nil
}

func (it *PubSub) Subscribe(topic string) <-chan int {
	logw.Infof("^PubSub.Subscribe - topic: %s", topic)
	defer logw.Infof("$PubSub.Subcribe")

	it.mu.Lock()
	defer it.mu.Unlock()
	logw.Debug("PubSub.Subscribe - acquired mutex")
	defer logw.Debug("PubSub.Subscribe - releasing mutex")

	sub := make(chan int, it.bufferSize)
	it.submap[topic] = append(it.submap[topic], sub)

	return sub
}

func (it *PubSub) Close() {
	logw.Debugf("^PubSub.Close")
	defer logw.Debugf("$PubSub.Close")

	if it.closed {
		return
	}

	it.mu.Lock()
	defer it.mu.Unlock()
	logw.Debug("PubSub.Close - acquired mutex")
	defer logw.Debug("PubSub.Close - releasing mutex")

	it.closed = true
	for _, topic := range it.submap {
		for _, sub := range topic {
			close(sub)
		}
	}
}

func (it *PubSub) Init(bufferSize int) {
	logw.Debug("^$PubSub.Init")
	it.submap = map[string][]chan int{}
	it.bufferSize = bufferSize
}

var Broker IPubSub = &PubSub{}
