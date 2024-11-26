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
}

type Register struct {
	Name    string
	Inputs  []string
	Outputs []string
	RunFunc     func(ctx context.Context, it IRegister)
	running bool
}

func (it *Register) GetName() string {
    logw.Debug("^$Register.GetName")
	return it.Name
}

func (it *Register) GetInputs() []string {
    logw.Debug("^$Register.GetInputs")
	return it.Inputs
}

func (it *Register) GetOutputs() []string {
    logw.Debug("^$Register.GetOutputs")
	return it.Outputs
}

func (it *Register) Run(ctx context.Context) {
    logw.Debug("^Register.Run")
    defer logw.Debug("$Register.Run")

	if it.running {
		return
	}

	go it.RunFunc(ctx, it)
	it.running = true
}

func RegisterRunDef(ctx context.Context, it IRegister) {
    name := it.GetName()
    logw.Debugf("^Register.RegisterRunDef: %s", name)
    defer logw.Debugf("$Register.RegisterRunDef: %s", name)

    oe := Broker.Subscribe(fmt.Sprintf("%s_OE", name))
    we := Broker.Subscribe(fmt.Sprintf("%s_WE", name))
    di := Broker.Subscribe("D")
    clk := Broker.Subscribe("CLK")
    rst := Broker.Subscribe("RST")

    for {
        select {
        case <-ctx.Done():
            return
        case dat := <-oe:
            logw.Infof("Register %s received OE update %s", name, dat)
        case dat := <-we:
            logw.Infof("Register %s received OE update %s", name, dat)
        case dat := <-di:
            logw.Infof("Register %s received D update %s", name, dat)
        case dat := <-clk:
            logw.Infof("Register %s received CLK update %s", name, dat)
        case dat := <-rst:
            logw.Infof("Register %s received RST update %s", name, dat)
        }
    }
}
