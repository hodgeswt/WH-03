package wh03

import (
	"github.com/hodgeswt/WH-03/internal/types"
)

var regA types.IRegister = &types.Register{
	Name:    "RegisterA",
	Inputs:  []string{"RST", "CLK", "D", "OE", "WE"},
	Outputs: []string{"DO"},
	RunFunc: types.RegisterRunDef,
}
