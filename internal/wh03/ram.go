package wh03

import (
	"context"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type IRam interface {
	ICircuit
}

type Ram struct {
	Size   int
	buffer map[string]int
	data   map[int]int
}

func (it *Ram) Buffer(key string, data int) {
	logw.Debug("^Ram.Buffer")
	defer logw.Debug("$Ram.Buffer")

	if it.buffer == nil {
		it.buffer = map[string]int{}
	}

	it.buffer[key] = data
}

func (it *Ram) Run(ctx context.Context) {
	logw.Debug("^Ram.Run")
	defer logw.Debug("$Ram.Run")

	oe := Broker.Subscribe("Ram_OE")
	we := Broker.Subscribe("Ram_WE")
	d := Broker.Subscribe("D")
	clk := Broker.Subscribe("CLK")
	rst := Broker.Subscribe("RST")
	memAdd := Broker.Subscribe("MEM_ADD")

	it.data = map[int]int{}

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-oe:
			logw.Infof("RAM received OE: %08b", dat)
			it.Buffer("OE", dat)
		case dat := <-we:
			logw.Infof("RAM received WE: %08b", dat)
			it.Buffer("WE", dat)
		case dat := <-d:
			it.Buffer("D", dat)
		case dat := <-memAdd:
			it.Buffer("MEM_ADD", dat)
		case dat := <-clk:
			if dat == 0 {
				// Falling edge
				it.UpdateState()
			}
		case dat := <-rst:
			logw.Infof("Ram received RST: %08b", dat)
			it.Reset()
		}
	}
}

func (it *Ram) Reset() {
	logw.Debug("^$Ram.Reset()")
	it.buffer = map[string]int{}
	it.data = map[int]int{}
}

func (it *Ram) UpdateState() {
	logw.Debug("^Ram.UpdateState")
	defer logw.Debug("$Ram.UpdateState")

	if it.buffer["WE"] == 1 {
		it.data[it.buffer["MEM_ADD"]] = it.buffer["D"]
	}

	toOut := it.data[it.buffer["MEM_ADD"]]

	if it.buffer["OE"] == 1 {
		Broker.Publish("D", toOut)
	}

    Broker.Publish("Ram_D", toOut)

}

var ram IRam = &Ram{}
