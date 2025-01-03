package wh03

import (
	"context"
	"fmt"

	"github.com/hodgeswt/WH-03/internal/util"
	"github.com/hodgeswt/utilw/pkg/logw"
)

var regA IRegister = &Register{
	Name:            "A",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
	OutToCustom:     true,
}

var regB IRegister = &Register{
	Name:            "B",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
	OutToCustom:     true,
}

var regC IRegister = &Register{
	Name:            "C",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
}

var output1 IRegister = &Register{
	Name:            "Output1",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
}

var output2 IRegister = &Register{
	Name:            "Output2",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
}

var acc IRegister = &Register{
	Name:            "Accumulator",
	Inputs:          []string{"RST", "CLK", "Alu_D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         AccRun,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
}

var mar IRegister = &Register{
	Name:            "MemoryAddress",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: MemoryAddressUpdateState,
	ResetFunc:       RegisterResetDef,
}

var instr IRegister = &Register{
	Name:            "Instruction",
	Inputs:          []string{"RST", "CLK", "Ram_D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
}

var flags IRegister = &Register{
	Name:            "Flags",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         FlagsRegisterRun,
	UpdateStateFunc: FlagsRegisterUpdateState,
	ResetFunc:       RegisterResetDef,
}

func MemoryAddressUpdateState(it *Register) {
	logw.Debugf("^registers.MemoryAddressUpdateState - %s", it.Name)
	defer logw.Debugf("$registers.MemoryAddressUpdateState - %s", it.Name)

	if it.InputBuffer["WE"] == 1 {
		it.State = it.InputBuffer["D"]
	}

	if it.InputBuffer["OE"] == 1 {
		Broker.Publish("D", it.State)
	}

	Broker.Publish("Mar_D", it.State)

	it.InputBuffer = map[string]int{}
}


func InstRegisterUpdateState(it *Register) {
	logw.Debugf("^registers.InstRegisterUpdateState - %s", it.Name)
	defer logw.Debugf("$registers.InstRegisterUpdateState - %s", it.Name)

	if it.InputBuffer["WE"] == 1 {
		it.State = it.InputBuffer["D"]
	}

	if it.InputBuffer["OE"] == 1 {
		Broker.Publish("D", it.State)
	}

	Broker.Publish("INST", it.State)

	it.InputBuffer = map[string]int{}
}

func FlagsRegisterUpdateState(it *Register) {
	logw.Debugf("^registers.FlagsRegisterUpdateState - %s", it.Name)
	defer logw.Debugf("$registers.FlagsRegisterUpdateState- %s", it.Name)

	if it.InputBuffer["WE"] == 1 {
		if it.InputBuffer["FD"] > util.MaxBit {
			panic("Flags register received invalid data")
		}

		it.State = it.InputBuffer["FD"]
	}

	Broker.Publish("Flags_GT", util.GetBit(it.State, 0))
	Broker.Publish("Flags_LT", util.GetBit(it.State, 1))
	Broker.Publish("Flags_EQ", util.GetBit(it.State, 2))
	Broker.Publish("Flags_C", util.GetBit(it.State, 3))

	it.InputBuffer = map[string]int{}
}

func AccRun(ctx context.Context, it *Register) {
	logw.Debugf("^Register.RegisterRunDef: %s", it.Name)
	defer logw.Debugf("$Register.RegisterRunDef: %s", it.Name)

	oe := Broker.Subscribe(fmt.Sprintf("%s_OE", it.Name))
	we := Broker.Subscribe(fmt.Sprintf("%s_WE", it.Name))

	d := Broker.Subscribe("Alu_D")
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
			logw.Infof("Register %s received WE update %08b", it.Name, dat)
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

func FlagsRegisterRun(ctx context.Context, it *Register) {
	logw.Debugf("^registers.FlagsRegisterRun: %s", it.Name)
	defer logw.Debugf("$registers.FlagsRegisterRun: %s", it.Name)

	fd := Broker.Subscribe("FD")
	clk := Broker.Subscribe("CLK")
	rst := Broker.Subscribe("RST")

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-fd:
			it.Buffer("FD", dat)
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
