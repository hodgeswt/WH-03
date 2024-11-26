package types

import (
	"context"
	"fmt"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type IRegister interface {
	ICircuit
	GetName() string
	GetInputs() []string
	GetOutputs() []string
	Buffer(key string, data string)
	UpdateState()
	Reset()
}

type Register struct {
	Name            string
	Inputs          []string
	Outputs         []string
	RunFunc         func(ctx context.Context, it *Register)
	UpdateStateFunc func(it *Register)
	ResetFunc       func(it *Register)
	running         bool
	buffer          map[string]string
	state           string
}

func (it *Register) GetName() string {
	logw.Debug("^$Register.GetName")
	return it.Name
}

func (it *Register) GetInputs() []string {
	logw.Debugf("^$Register.GetInputs - %s", it.Name)
	return it.Inputs
}

func (it *Register) GetOutputs() []string {
	logw.Debugf("^$Register.GetOutputs - %s", it.Name)
	return it.Outputs
}

func (it *Register) Run(ctx context.Context) {
	logw.Debugf("^Register.Run - %s", it.Name)
	defer logw.Debugf("$Register.Run - %s", it.Name)

	if it.running {
		return
	}

	go it.RunFunc(ctx, it)
	it.running = true
}

func (it *Register) UpdateState() {
	logw.Debugf("^Register.UpdateState - %s", it.Name)
	defer logw.Debugf("$Register.UpdateState - %s", it.Name)

	it.UpdateStateFunc(it)
}

func (it *Register) Buffer(key string, data string) {
	logw.Debugf("^$Register.Buffer - %s", it.Name)
	if it.buffer == nil {
		it.buffer = map[string]string{}
	}
	it.buffer[key] = data
}

func (it *Register) Reset() {
	logw.Debugf("^$Register.Reset - %s", it.Name)

	it.ResetFunc(it)
}

func RegisterUpdateStateDef(it *Register) {
	logw.Debugf("^Register.RegisterUpdateStateDef - %s", it.Name)
	defer logw.Debugf("$Register.RegisterUpdateStateDef - %s", it.Name)

	if it.buffer["WE"] == "1" && it.buffer["D"] != "" {
		it.state = it.buffer["D"]
	}

	if it.buffer["OE"] == "1" && it.state != "" {
		Broker.Publish("D", it.state)
	}

	it.buffer = map[string]string{}
}

func RegisterResetDef(it *Register) {
	logw.Debugf("^$Register.RegisterResetDef: %s", it.Name)

	clear(it.buffer)
	it.state = "00000000"
}

func RegisterRunDef(ctx context.Context, it *Register) {
	logw.Debugf("^Register.RegisterRunDef: %s", it.Name)
	defer logw.Debugf("$Register.RegisterRunDef: %s", it.Name)

	oe := Broker.Subscribe(fmt.Sprintf("%s_OE", it.Name))
	we := Broker.Subscribe(fmt.Sprintf("%s_WE", it.Name))
	d := Broker.Subscribe("D")
	clk := Broker.Subscribe("CLK")
	rst := Broker.Subscribe("RST")

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-oe:
			logw.Infof("Register %s received OE update %s", it.Name, dat)
			it.Buffer("OE", dat)
		case dat := <-we:
			logw.Infof("Register %s received OE update %s", it.Name, dat)
			it.Buffer("WE", dat)
		case dat := <-d:
			it.Buffer("D", dat)
		case dat := <-clk:
			if dat == "1" {
				// Rising edge
				it.UpdateState()
			}
		case dat := <-rst:
			logw.Infof("Register %s received RST update %s", it.Name, dat)
			it.Reset()
		}
	}
}
