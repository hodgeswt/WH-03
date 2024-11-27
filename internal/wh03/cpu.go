package wh03

import (
	"context"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type CPU struct {
	ctx       context.Context
	cancel    context.CancelFunc
	registers []IRegister
	clock     ICircuit
	prgc      ICircuit
	stepctr   IStepCounter
	ram       IRam
	Cfg       *Config
}

func (it *CPU) Setup() {
	logw.Debug("^wh03.Setup")
	defer logw.Debug("$wh03.Setup")
	it.registers = []IRegister{
		regA,
		regB,
		regC,
		output1,
		output2,
		acc,
		flags,
		mar,
		instr,
	}

	it.clock = &Clock{
		Freq: it.Cfg.ClockFreq,
	}
	it.stepctr = &StepCounter{Limit: 7}
	it.prgc = &ProgramCounter{}
	it.ram = &Ram{Size: it.Cfg.RamK * (2 ^ 10)}

	it.ctx, it.cancel = context.WithCancel(context.Background())
}

func (it *CPU) Run() {
	logw.Debug("^wh03.Run")
	defer logw.Debug("$wh03.Run")

	it.Setup()

	go it.clock.Run(it.ctx)
	go it.prgc.Run(it.ctx)
	go it.stepctr.Run(it.ctx)
	go it.ram.Run(it.ctx)

	for _, register := range it.registers {
		go register.Run(it.ctx)
	}

	it.run()
}

func (it *CPU) run() {
	clk := Broker.Subscribe("CLK")
	d := []map[string]string{
		{
			"D":    "01010101",
			"A_WE": "1",
		},
		{
			"A_OE": "1",
		},
		{
			"D": "11111111",
		},
	}
	i := 0
	for {
		select {
		case <-it.ctx.Done():
			logw.Debug("wh03.run - context canceled")
			return
		case dat := <-clk:
			if dat == "1" {
				if i <= len(d)-1 {
					inst := d[i]
					i = i + 1
					for k, v := range inst {
						Broker.Publish(k, v)
					}
				}
			}
		}
	}
}

func (it *CPU) Stop() {
	logw.Debug("^wh03.Stop")
	defer logw.Debug("$wh03.Stop")
	it.cancel()
}
