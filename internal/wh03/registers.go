package wh03

import (
	"github.com/hodgeswt/WH-03/internal/types"
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

var flags types.IRegister = &types.Register{
	Name:            "Flags",
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
