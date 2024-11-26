package wh03

import (
	"context"

	"github.com/hodgeswt/WH-03/internal/types"
	"github.com/hodgeswt/utilw/pkg/logw"
)

var regA types.IRegister = &types.Register{
	Name:            "A",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         types.RegisterRunDef,
	UpdateStateFunc: types.RegisterUpdateStateDef,
	ResetFunc:       types.RegisterResetDef,
}

var regB types.IRegister = &types.Register{
	Name:            "B",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         types.RegisterRunDef,
	UpdateStateFunc: types.RegisterUpdateStateDef,
	ResetFunc:       types.RegisterResetDef,
}

var regC types.IRegister = &types.Register{
	Name:            "C",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         types.RegisterRunDef,
	UpdateStateFunc: types.RegisterUpdateStateDef,
	ResetFunc:       types.RegisterResetDef,
}

var output1 types.IRegister = &types.Register{
	Name:            "Output1",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         types.RegisterRunDef,
	UpdateStateFunc: types.RegisterUpdateStateDef,
	ResetFunc:       types.RegisterResetDef,
}

var output2 types.IRegister = &types.Register{
	Name:            "Output2",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         types.RegisterRunDef,
	UpdateStateFunc: types.RegisterUpdateStateDef,
	ResetFunc:       types.RegisterResetDef,
}

var acc types.IRegister = &types.Register{
	Name:            "Accumulator",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         types.RegisterRunDef,
	UpdateStateFunc: types.RegisterUpdateStateDef,
	ResetFunc:       types.RegisterResetDef,
}

var mar types.IRegister = &types.Register{
	Name:            "MemoryAddress",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         types.RegisterRunDef,
	UpdateStateFunc: types.RegisterUpdateStateDef,
	ResetFunc:       types.RegisterResetDef,
}

var instr types.IRegister = &types.Register{
	Name:            "Instruction",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         types.RegisterRunDef,
	UpdateStateFunc: types.RegisterUpdateStateDef,
	ResetFunc:       types.RegisterResetDef,
}

var flags types.IRegister = &types.Register{
	Name:            "Flags",
	Inputs:          []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs:         []string{"DO"},
	RunFunc:         FlagsRegisterRun,
	UpdateStateFunc: FlagsRegisterUpdateState,
	ResetFunc:       types.RegisterResetDef,
}

func FlagsRegisterUpdateState(it *types.Register) {
	logw.Debugf("^registers.FlagsRegisterUpdateState - %s", it.Name)
	defer logw.Debugf("$registers.FlagsRegisterUpdateState- %s", it.Name)

	if it.InputBuffer["WE"] == "1" && it.InputBuffer["FD"] != "" {
		if len(it.InputBuffer["FD"]) != 8 {
			panic("Flags register received invalid data")
		}

		it.State = it.InputBuffer["FD"]
	}

	types.Broker.Publish("Flags_GT", string(it.State[0]))
	types.Broker.Publish("Flags_EQ", string(it.State[1]))
	types.Broker.Publish("Flags_LT", string(it.State[2]))
	types.Broker.Publish("Flags_C", string(it.State[3]))

	it.InputBuffer = map[string]string{}
}

func FlagsRegisterRun(ctx context.Context, it *types.Register) {
	logw.Debugf("^registers.FlagsRegisterRun: %s", it.Name)
	defer logw.Debugf("$registers.FlagsRegisterRun: %s", it.Name)

	fd := types.Broker.Subscribe("FD")
	clk := types.Broker.Subscribe("CLK")
	rst := types.Broker.Subscribe("RST")

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
