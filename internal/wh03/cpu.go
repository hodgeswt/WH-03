package wh03

import (
	"context"

	"github.com/hodgeswt/WH-03/internal/types"
	"github.com/hodgeswt/utilw/pkg/logw"
)

type CPU struct {
	ctx    context.Context
	cancel context.CancelFunc
	RegA   types.IRegister
    Clock  types.ICircuit
    Cfg    *Config
}

type Config struct {
    ClockFreq int
}

func (it *CPU) Setup() {
	logw.Debug("^wh03.Setup")
	defer logw.Debug("$wh03.Setup")
	it.RegA = regA
    it.Clock = &types.Clock{
        Freq: it.Cfg.ClockFreq,
    }
	it.ctx, it.cancel = context.WithCancel(context.Background())
}

func (it *CPU) Run() {
	logw.Debug("^wh03.Run")
	defer logw.Debug("$wh03.Run")

	it.Setup()

    go it.Clock.Run(it.ctx)
	go it.RegA.Run(it.ctx)
	it.run()
}

func (it *CPU) run() {
	for {
		select {
		case <-it.ctx.Done():
			logw.Debug("wh03.run - context canceled")
			return
		}
	}
}

func (it *CPU) Stop() {
	logw.Debug("^wh03.Stop")
	defer logw.Debug("$wh03.Stop")
	it.cancel()
}

