package wh03

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hodgeswt/WH-03/internal/util"
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
	rom       IRom
	ctrl      ICtrl
	alu       IAlu
	stack     IStack
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
	it.stepctr = &StepCounter{Limit: it.Cfg.StepsPerInstruction}
	it.prgc = &ProgramCounter{}
	it.ram = &Ram{Size: it.Cfg.RamK * util.IntPow(2, 10)}
	it.rom = &Rom{Size: it.Cfg.RomK * util.IntPow(2, 10)}
	it.ctrl = &Ctrl{}
	it.alu = &Alu{}

	size, err := strconv.ParseInt(it.Cfg.StackStart, 16, 64)
	if err != nil {
		panic(fmt.Sprintf("Invalid hex value provided for stack start: %s", it.Cfg.StackStart))
	}
	it.stack = &Stack{StackStart: int(size)}

	err = it.rom.Load(it.Cfg.RomFile)

	if err != nil {
		panic(fmt.Sprintf("Unable to load ROM file at %s, err: %v", it.Cfg.RomFile, err))
	}

	it.ctx, it.cancel = context.WithCancel(context.Background())
}

func (it *CPU) Run() {
	logw.Debug("^wh03.Run")
	defer logw.Debug("$wh03.Run")

	it.Setup()

	go it.rom.Run(it.ctx)
	go it.clock.Run(it.ctx)
	go it.prgc.Run(it.ctx)
	go it.stepctr.Run(it.ctx)
	go it.ram.Run(it.ctx)
	go it.ctrl.Run(it.ctx)
	go it.alu.Run(it.ctx)
	go it.stack.Run(it.ctx)

	for _, register := range it.registers {
		go register.Run(it.ctx)
	}

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
