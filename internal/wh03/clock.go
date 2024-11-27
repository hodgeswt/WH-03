package wh03

import (
	"context"
	"fmt"
	"time"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type Clock struct {
	Freq  int
	state int
}

func (it *Clock) Run(ctx context.Context) {
	freq := time.Duration(it.Freq)
	it.state = 0
	for {
		select {
		case <-ctx.Done():
			logw.Debugf("Clock.Run - context cancelled")
			return
		case <-time.After(freq * time.Second):
			Broker.Publish("CLK", fmt.Sprintf("%d", it.state))
			toggle(&it.state)
		}
	}
}

func (it *Clock) Reset() {
	it.state = 0
}

func (it *Clock) Buffer(key string, data string) {}
func (it *Clock) UpdateState()                   {}

func toggle(int *int) {
	if *int == 0 {
		*int = 1
	} else {
		*int = 0
	}
}

