package wh03

import (
	"context"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type IStepCounter interface {
	ICircuit
}

type StepCounter struct {
	buffer map[string]string
	state  int
    Limit  int
}

func (it *StepCounter) Run(ctx context.Context) {
	logw.Debugf("^StepCounter.Run")
	defer logw.Debugf("$StepCounter.Run")

	clk := Broker.Subscribe("CLK")
	rst := Broker.Subscribe("RST")
    stepRst := Broker.Subscribe("StepCounter_RST")

	it.state = 0

	for {
		select {
		case <-ctx.Done():
			logw.Debugf("StepCounter.Run - context cancelled")
			return
		case dat := <-clk:
			if dat == "0" {
				// Falling edge
				it.UpdateState()
			}
        case dat := <-stepRst:
            logw.Infof("Step Counter received StepCounter_RST: %s", dat)
            it.Reset()
		case dat := <-rst:
			logw.Infof("Step Counter received RST update %s", dat)
			it.Reset()
		}
	}

}

func (it *StepCounter) Buffer(key string, data string) {
	logw.Debugf("^StepCounter.Buffer")
	defer logw.Debugf("$StepCounter.Buffer")

	if it.buffer == nil {
		it.buffer = map[string]string{}
	}

	it.buffer[key] = data
}

func (it *StepCounter) UpdateState() {
	logw.Debugf("^StepCounter.UpdateState")
	defer logw.Debugf("$StepCounter.UpdateState")

    it.state = it.state + 1

    if it.state >= it.Limit {
        it.state = 0
    }

}

func (it *StepCounter) Reset() {
	logw.Debugf("^$StepCounter.Reset")
    it.state = 0
}

var stpctr IStepCounter = &StepCounter{}
