package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/hodgeswt/WH-03/internal/util"
)

type Exception struct {
	Key     string   `json:"key"`
	When    []string `json:"when"`
	Formula []string `json:"formula"`
}
type Formula struct {
	Inputs     []string    `json:"inputs"`
	Formula    []string    `json:"formula"`
	Name       string      `json:"name"`
	Exceptions []Exception `json:"exceptions"`
}

type Instruction struct {
	Name  string   `json:"name"`
	Steps []string `json:"steps"`
}

type RomDef struct {
	SizeK               int               `json:"sizeK"`
	RomFile             string            `json:"romFile"`
	LabelFile           string            `json:"labelFile"`
	StepsPerInstruction int               `json:"stepsPerInstruction"`
	Definitions         map[string]string `json:"definitions"`
	Base                []string          `json:"base"`
	Formulae            []Formula         `json:"formulae"`
	Instructions        []Instruction     `json:"instructions"`
}

func hexToInt(data string) int64 {
	conv, err := strconv.ParseInt(data, 16, 64)

	if err != nil {
		panic(fmt.Sprintf("Unable to convert %s to hex: %v", data, err))
	}

	return conv
}

func stepToHex(step string, defs map[string]string) int64 {
	fmt.Printf("Converting %s\n", step)
	steps := strings.Split(step, "+")
	fmt.Printf("Converting %v\n", step)
	var x int64 = 0

	for _, s := range steps {
		x = x | hexToInt(defs[s])
	}

	return x
}

func addBase(addr int64, base []string, defs map[string]string, rom *[]int64) {
	var i int64 = 0
	for _, step := range base {
		(*rom)[(addr | i)] = stepToHex(step, defs)
		i++
	}
}

func fillEnd(i int64, maximum int, addr int64, defs map[string]string, rom *[]int64) {
	setReset := false
	for {
		if i >= int64(maximum) {
			break
		}

		if !setReset {
			(*rom)[(addr | i)] = hexToInt(defs["reset_stepc"])
		} else {
			(*rom)[(addr | i)] = hexToInt(defs["nop"])
		}

		i++
	}

}

func main() {
	// TODO: all logic
	content, err := os.ReadFile("rom.json")

	var cfg *RomDef
	if err != nil {
		panic("Unable to find rom.json.")
	} else {
		loadedCfg := new(RomDef)
		err := json.Unmarshal(content, &loadedCfg)
		if err != nil {
			panic("Unable to load rom.json.")
		} else {
			cfg = loadedCfg
		}
	}

	fmt.Printf("%+v\n", cfg)

	fmt.Printf("Making ROM of size %d\n", cfg.SizeK*util.IntPow(2, 10))
	rom := make([]int64, cfg.SizeK*util.IntPow(2, 10))
	labels := map[int]string{}

	label := 0

	for _, inst := range cfg.Instructions {
		if len(inst.Steps) >= (cfg.StepsPerInstruction - len(cfg.Base)) {
			panic(fmt.Sprintf("Instruction %v is too long.", inst))
		}
		addr := int64(label) << 8

		labels[label] = inst.Name
		label++

		var i int64 = int64(len(cfg.Base) - 1)

		addBase(addr, cfg.Base, cfg.Definitions, &rom)

		for _, step := range inst.Steps {
			rom[(addr | i)] = stepToHex(step, cfg.Definitions)
			i++
		}

		fillEnd(i, cfg.StepsPerInstruction, addr, cfg.Definitions, &rom)
	}

	for _, formula := range cfg.Formulae {
		for _, x := range formula.Inputs {
			for _, y := range formula.Inputs {
				if x == y {
					continue
				}

				name := strings.ReplaceAll(formula.Name, "@1", x)
				name = strings.ReplaceAll(name, "@2", y)
				addr := int64(label) << 8
				labels[label] = name
				label++

				doException := false
				var exception Exception

				for _, excp := range formula.Exceptions {
					v := ""
					if excp.Key == "@1" {
						v = x
					} else {
						v = y
					}

					if slices.Contains(excp.When, v) {
						doException = true
						exception = excp
						break
					}
				}

				addBase(addr, cfg.Base, cfg.Definitions, &rom)
				var i int64 = int64(len(cfg.Base) - 1)

				var formulaSteps []string

				if doException {
					formulaSteps = exception.Formula
				} else {
					formulaSteps = formula.Formula
				}

				for _, step := range formulaSteps {
					f := strings.ReplaceAll(step, "@1", x)
					f = strings.ReplaceAll(f, "@2", y)

					rom[(addr | i)] = stepToHex(f, cfg.Definitions)
					i++
				}

				fillEnd(i, cfg.StepsPerInstruction, addr, cfg.Definitions, &rom)
			}
		}
	}

	file, err := os.Create(cfg.RomFile)
	if err != nil {
		panic(fmt.Sprintf("Unable to create ROM file at %s, err: %v", cfg.RomFile, err))
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, rom)

	if err != nil {
		panic(fmt.Sprintf("Unable to write ROM file at %s, err: %v", cfg.RomFile, err))
	}

	labelFile, err := os.Create(cfg.LabelFile)
	if err != nil {
		panic(fmt.Sprintf("Unable to create Label file at %s, err: %v", cfg.RomFile, err))
	}
	defer labelFile.Close()

    for k, v := range labels {
        fmt.Printf("Writing label: %04x\n", k)
        labelFile.WriteString(fmt.Sprintf("%04x: %s\n", k, v))
    }
}
