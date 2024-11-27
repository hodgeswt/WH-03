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
	buffer map[string]string
	data   map[string]string
}

func (it *Ram) Buffer(key string, data string) {
	logw.Debug("^Ram.Buffer")
	defer logw.Debug("$Ram.Buffer")

	if it.buffer == nil {
		it.buffer = map[string]string{}
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

	it.data = map[string]string{}

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-oe:
			logw.Infof("RAM received OE: %s", dat)
			it.Buffer("OE", dat)
		case dat := <-we:
			logw.Infof("RAM received WE: %s", dat)
			it.Buffer("WE", dat)
		case dat := <-d:
			it.Buffer("D", dat)
		case dat := <-memAdd:
			it.Buffer("MEM_ADD", dat)
		case dat := <-clk:
			if dat == "0" {
				// Falling edge
				it.UpdateState()
			}
		case dat := <-rst:
			logw.Infof("Ram received RST: %s", dat)
			it.Reset()
		}
	}
}

func (it *Ram) Reset() {
	logw.Debug("^$Ram.Reset()")
	it.buffer = map[string]string{}
	it.data = map[string]string{}
}

func (it *Ram) UpdateState() {
	logw.Debug("^Ram.UpdateState")
	defer logw.Debug("$Ram.UpdateState")

	if it.buffer["WE"] == "1" && it.buffer["MEM_ADD"] != "" {
		toWrite := it.buffer["D"]
		if toWrite == "" {
			toWrite = "00000000"
		}

		it.data[it.buffer["MEM_ADD"]] = toWrite
	}

	toOut := it.data[it.buffer["MEM_ADD"]]
	if toOut == "" {
		toOut = "00000000"
	}

	if it.buffer["OE"] == "1" && it.buffer["MEM_ADD"] != "" {
		Broker.Publish("D", toOut)
	}

    Broker.Publish("Ram_D", toOut)

}

var ram IRam = &Ram{}
