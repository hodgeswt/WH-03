package wh03

import (
	"context"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type IProgramCounter interface {
	ICircuit
	Increment()
	Reset()
	Set(val int)
	Buffer(key string, data int)
}

type ProgramCounter struct {
	count  int
	buffer map[string]int
}

func (it *ProgramCounter) Buffer(key string, data int) {
	logw.Debug("^ProgramCounter.Buffer")
	defer logw.Debug("$ProgramCounter.Buffer")

	if it.buffer == nil {
		it.buffer = map[string]int{}
	}

	it.buffer[key] = data
}

func (it *ProgramCounter) Run(ctx context.Context) {
	logw.Debug("^ProgramCounter.Run")
	defer logw.Debug("$ProgramCounter.Run")

	we := Broker.Subscribe("PRGC_WE")
	clk := Broker.Subscribe("CLK")
	d := Broker.Subscribe("D")
	rst := Broker.Subscribe("RST")
	e := Broker.Subscribe("PRGC_E")
	oe := Broker.Subscribe("PRGC_OE")

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-we:
			logw.Infof("ProgramCounter received OE update %08d", dat)
			it.Buffer("WE", dat)
		case dat := <-d:
			it.Buffer("D", dat)
		case dat := <-clk:
			if dat == 1 {
				// Rising edge
				it.UpdateState()
			}
		case dat := <-e:
			it.Buffer("E", dat)
		case dat := <-oe:
			it.Buffer("OE", dat)
		case dat := <-rst:
			logw.Infof("ProgramCounter received RST update %08d", dat)
			it.Reset()
		}
	}

}

func (it *ProgramCounter) Increment() {
	logw.Debug("^$ProgramCounter.Increment")
	it.count = it.count + 1
}

func (it *ProgramCounter) Reset() {
	logw.Debug("^$ProgramCounter.Reset")
	it.count = 0
}

func (it *ProgramCounter) Set(val int) {
	logw.Debug("^$ProgramCounter.Set")
	it.count = val
}

func (it *ProgramCounter) UpdateState() {
	logw.Debug("^ProgramCounter.UpdateState")
	defer logw.Debug("$ProgramCounter.UpdateState")

	if it.buffer == nil || len(it.buffer) == 0 {
		return
	}

	if it.buffer["E"] == 1 {
		it.Increment()
	}

	if it.buffer["WE"] == 1 {
		it.count = it.buffer["D"]
	}

	if it.buffer["OE"] == 1 {
		Broker.Publish("D", it.count)
	}

	it.buffer = map[string]int{}
}
