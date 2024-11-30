//go:build unit
// +build unit

package wh03

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClock_Reset(t *testing.T) {
	clk := &Clock{state: 10}

	clk.Reset()

	assert.Equal(t, 0, clk.state)
}

func TestClock_Run(t *testing.T) {
	clk := &Clock{Freq: 1}
	Broker = &PubSub{}
	Broker.Init(10, false)

	ctx, cancel := context.WithCancel(context.Background())

	ch := Broker.Subscribe("CLK")

	go clk.Run(ctx)

	expected := 0
loop:
	for {
		select {
		case d := <-ch:
			assert.Equal(t, expected, d)
			if expected == 1 {
				break loop
			}
			expected = 1
		case <-time.After(4 * time.Second):
			assert.FailNow(t, "Did not receive tick in allotted time")
		}
	}

	cancel()
}

func TestClock_Halt(t *testing.T) {
    clk := &Clock{Freq: 1, testMode: true, halt: true}
	Broker = &PubSub{}
	Broker.Init(10, false)

	ctx, cancel := context.WithCancel(context.Background())

	ch := Broker.Subscribe("CLK")

	go clk.Run(ctx)

loop:
	for {
		select {
		case <-ch:
			assert.FailNow(t, "Received unexpected tick")
		case <-time.After(2 * time.Second):
			break loop
		}
	}

	cancel()
}

func TestToggle_Zero(t *testing.T) {
	x := 0

	toggle(&x)

	assert.Equal(t, 1, x)
}

func TestToggle_One(t *testing.T) {
	x := 1

	toggle(&x)

	assert.Equal(t, 0, x)
}

func TestToggle_Other(t *testing.T) {
	x := 3

	toggle(&x)

	assert.Equal(t, 0, x)
}
