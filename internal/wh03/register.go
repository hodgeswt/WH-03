package wh03

import (
	"context"
	"fmt"
	"slices"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type IRegister interface {
	ICircuit
	GetName() string
	GetInputs() []string
	GetOutputs() []string
	Buffer(key string, data int)
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
	InputBuffer     map[string]int
	State           int
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

func (it *Register) Buffer(key string, data int) {
	logw.Debugf("^$Register.Buffer - %s", it.Name)
	if it.InputBuffer == nil {
		it.InputBuffer = map[string]int{}
	}
	it.InputBuffer[key] = data
}

func (it *Register) Reset() {
	logw.Debugf("^$Register.Reset - %s", it.Name)

	it.ResetFunc(it)
}

func RegisterUpdateStateDef(it *Register) {
	logw.Debugf("^Register.RegisterUpdateStateDef - %s", it.Name)
	defer logw.Debugf("$Register.RegisterUpdateStateDef - %s", it.Name)

	if it.InputBuffer["WE"] == 1 {
		it.State = it.InputBuffer["D"]
	}

	if it.InputBuffer["OE"] == 1 {
		Broker.Publish("D", it.State)
	}

	it.InputBuffer = map[string]int{}
}

func RegisterResetDef(it *Register) {
	logw.Debugf("^$Register.RegisterResetDef: %s", it.Name)

	clear(it.InputBuffer)
	it.State = 0
}

func RegisterRunDef(ctx context.Context, it *Register) {
	logw.Debugf("^Register.RegisterRunDef: %s", it.Name)
	defer logw.Debugf("$Register.RegisterRunDef: %s", it.Name)

	oe := Broker.Subscribe(fmt.Sprintf("%s_OE", it.Name))
	we := Broker.Subscribe(fmt.Sprintf("%s_WE", it.Name))

	var d <-chan int
	if slices.Contains(it.Inputs, "Ram_D") {
		d = Broker.Subscribe("Ram_D")
	} else {
		d = Broker.Subscribe("D")
	}

	clk := Broker.Subscribe("CLK")
	rst := Broker.Subscribe("RST")

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-oe:
			logw.Infof("Register %s received OE update %08b", it.Name, dat)
			it.Buffer("OE", dat)
		case dat := <-we:
            logw.Infof("Register %s received OE update %08b", it.Name, dat)
			it.Buffer("WE", dat)
		case dat := <-d:
			it.Buffer("D", dat)
		case dat := <-clk:
			if dat == 1 {
				// Rising edge
				it.UpdateState()
			}
		case dat := <-rst:
			logw.Infof("Register %s received RST update %08b", it.Name, dat)
			it.Reset()
		}
	}
}
