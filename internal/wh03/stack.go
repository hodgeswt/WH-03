package wh03

import (
	"context"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type IStack interface {
	ICircuit
}

type Stack struct {
	buffer     map[string]int
	StackStart int
	ptr        int
}

func (it *Stack) Run(ctx context.Context) {
    logw.Debug("^Stack.Run")
	defer logw.Debug("$Stack.Run")

    it.ptr = it.StackStart


	oe := Broker.Subscribe("Stack_OE")
    inc := Broker.Subscribe("Stack_INC")
    dec := Broker.Subscribe("Stack_DEC")
	clk := Broker.Subscribe("CLK")
	rst := Broker.Subscribe("RST")

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-oe:
			it.Buffer("Stack_OE", dat)
		case dat := <-inc:
			it.Buffer("Stack_INC", dat)
		case dat := <-dec:
			it.Buffer("Stack_DEC", dat)
		case dat := <-clk:
			if dat == 1 {
				// Rising edge
				it.UpdateState()
			}
		case dat := <-rst:
			logw.Infof("Stack received RST update %08b", dat)
			it.Reset()
		}
	}
}

func (it *Stack) Buffer(key string, data int) {
	logw.Debug("^Stack.Buffer")
	defer logw.Debug("$Stack.Buffer")

	if it.buffer == nil {
		it.buffer = map[string]int{}
	}

	it.buffer[key] = data
}

func (it *Stack) UpdateState() {
    logw.Debug("^Stack.UpdateState")
	defer logw.Debug("$Stack.UpdateState")

    if it.buffer["Stack_INC"] == 1 {
        it.ptr++
    } else if it.buffer["Stack_DEC"] == 1 {
        it.ptr--
    }

    if it.buffer["Stack_OE"] == 1 {
        Broker.Publish("D", it.ptr)
    }

}

func (it *Stack) Reset() {
	it.buffer = map[string]int{}
    it.ptr = it.StackStart
}
