package wh03

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hodgeswt/utilw/pkg/logw"
	mqtt "github.com/mochi-mqtt/server/v2"

	"github.com/mochi-mqtt/server/v2/hooks/auth"

	"github.com/mochi-mqtt/server/v2/listeners"
)

type IPubSub interface {
	Publish(topic string, msg int) error
	Subscribe(topic string) <-chan int
	Close()
    Run(ctx context.Context)
	Init(bufferSize int, probeEnabled bool)
}

type PubSub struct {
	mu         sync.RWMutex
	submap     map[string][]chan int
	closed     bool
	bufferSize int
	server     *mqtt.Server
	running    bool
}

func (it *PubSub) Publish(topic string, msg int) error {
	logw.Debug("^PubSub.Publish")
	defer logw.Debug("$PubSub.Publish")

	if it.server != nil && !it.running {
		panic(errors.New("Probe enabled but server not running"))
	}

	it.mu.RLock()
	defer it.mu.RUnlock()

	logw.Infof("PubSub.Publish - topic: %s, msg: %08b", topic, msg)

	t := it.submap[topic]

	for _, sub := range t {
		sub <- msg
	}

	if it.server != nil {
		it.server.Publish(topic, []byte(fmt.Sprintf("%08b", msg)), true, 0)
	}

	return nil
}

func (it *PubSub) Subscribe(topic string) <-chan int {
	logw.Infof("^PubSub.Subscribe - topic: %s", topic)
	defer logw.Infof("$PubSub.Subcribe")

	it.mu.Lock()
	defer it.mu.Unlock()

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

	it.closed = true
	for _, topic := range it.submap {
		for _, sub := range topic {
			close(sub)
		}
	}

	if it.server != nil {
		it.server.Close()
	}
}

func (it *PubSub) Init(bufferSize int, probeEnabled bool) {
	logw.Debug("^$PubSub.Init")
	it.submap = map[string][]chan int{}
	it.bufferSize = bufferSize

	if probeEnabled {
		s := mqtt.New(&mqtt.Options{
			InlineClient: true,
		})
		_ = s.AddHook(new(auth.AllowHook), nil)
		tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: ":1883"})
		err := s.AddListener(tcp)
		if err != nil {
			panic("Unable to start probe server")
		}

		it.server = s

	}
}

func (it *PubSub) Run(ctx context.Context) {
	it.running = true
	err := it.server.Serve()
	if err != nil {
		panic(err)
	}
}

var Broker IPubSub = &PubSub{}
