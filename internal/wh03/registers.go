package wh03

import (
	"context"

	"github.com/hodgeswt/utilw/pkg/logw"
)

var regA IRegister = &Register{
	Name:            "A",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
}

var regB IRegister = &Register{
	Name:            "B",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
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
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
	ResetFunc:       RegisterResetDef,
}

var mar IRegister = &Register{
	Name:            "MemoryAddress",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         RegisterRunDef,
	UpdateStateFunc: RegisterUpdateStateDef,
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

func FlagsRegisterUpdateState(it *Register) {
	logw.Debugf("^registers.FlagsRegisterUpdateState - %s", it.Name)
	defer logw.Debugf("$registers.FlagsRegisterUpdateState- %s", it.Name)

	if it.InputBuffer["WE"] == "1" && it.InputBuffer["FD"] != "" {
		if len(it.InputBuffer["FD"]) != 8 {
			panic("Flags register received invalid data")
		}

		it.State = it.InputBuffer["FD"]
	}

	Broker.Publish("Flags_GT", string(it.State[0]))
	Broker.Publish("Flags_EQ", string(it.State[1]))
	Broker.Publish("Flags_LT", string(it.State[2]))
	Broker.Publish("Flags_C", string(it.State[3]))

	it.InputBuffer = map[string]string{}
}

func FlagsRegisterRun(ctx context.Context, it *Register) {
	logw.Debugf("^registers.FlagsRegisterRun: %s", it.Name)
	defer logw.Debugf("$registers.FlagsRegisterRun: %s", it.Name)

	fd := Broker.Subscribe("FD")
	clk := Broker.Subscribe("CLK")
	rst := Broker.Subscribe("RST")

    it.State = "00000000"

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-fd:
			it.Buffer("FD", dat)
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
