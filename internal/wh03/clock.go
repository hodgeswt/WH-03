package wh03

import (
	"context"
	"time"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type Clock struct {
	Freq  int
	state int
	halt  bool
}

func (it *Clock) Run(ctx context.Context) {
	freq := time.Duration(it.Freq)
	it.state = 0
	it.halt = false

	hlt := Broker.Subscribe("HLT")
	for {
		select {
		case <-ctx.Done():
			logw.Debugf("Clock.Run - context cancelled")
			return
		case dat := <-hlt:
			it.halt = dat == 1
		case <-time.After(freq * time.Second):
			if !it.halt {
				Broker.Publish("CLK", it.state)
				toggle(&it.state)

			}
		}
	}
}

func (it *Clock) Reset() {
	it.state = 0
}

func (it *Clock) Buffer(key string, data int) {}
func (it *Clock) UpdateState()                {}

func toggle(int *int) {
	if *int == 0 {
		*int = 1
	} else {
		*int = 0
	}
}
