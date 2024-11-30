package wh03

import (
	"context"

	"github.com/hodgeswt/WH-03/internal/util"
	"github.com/hodgeswt/utilw/pkg/logw"
)

type ICtrl interface {
	ICircuit
}

type Ctrl struct {
	buffer  map[string]int
	decodeA IDecoder
	decodeB IDecoder
}

func (it *Ctrl) Buffer(key string, data int) {
	logw.Debug("^Ctrl.Buffer")
	defer logw.Debug("$Ctrl.Buffer")

	if it.buffer == nil {
		it.buffer = map[string]int{}
	}

	it.buffer[key] = data
}

func (it *Ctrl) Reset() {
	logw.Debug("^$Ctrl.Reset")
	it.buffer = map[string]int{}
}

func (it *Ctrl) Run(ctx context.Context) {
	logw.Debug("^Ctrl.Run")
	defer logw.Debug("^Ctrl.Run")

	it.decodeA = &Decoder{
		Bitwidth: 4,
		Outputs: map[int]string{
			0:  "A_OE",
			1:  "B_OE",
			2:  "C_OE",
			3:  "Output1_OE",
			4:  "Output2_OE",
			5:  "Accumulator_OE",
			6:  "MemoryAddress_OE",
			7:  "Instruction_OE",
			8:  "Flags_OE",
			9:  "Ram_OE",
			10: "PRGC_OE",
		},
	}

	it.decodeB = &Decoder{
		Bitwidth: 4,
		Outputs: map[int]string{
			0:  "A_WE",
			1:  "B_WE",
			2:  "C_WE",
			3:  "Output1_WE",
			4:  "Output2_WE",
			5:  "Accumulator_WE",
			6:  "MemoryAddress_WE",
			7:  "Instruction_WE",
			8:  "Flags_WE",
			9:  "Ram_WE",
			10: "PRGC_WE",
		},
	}

	d := Broker.Subscribe("Rom_D")
	rst := Broker.Subscribe("RST")

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-d:
			it.Buffer("Rom_D", dat)
			it.UpdateState()
		case dat := <-rst:
			logw.Infof("Rom received RST update %08b", dat)
			it.Reset()
		}
	}
}

func (it *Ctrl) UpdateState() {
	logw.Debug("^Ctrl.Run")
	defer logw.Debug("^Ctrl.Run")

	d := it.buffer["Rom_D"]
	logw.Debugf("Ctrl.Run - Rom_D: %04x", d)
	enableDecodeA := util.GetBit(d, 4) == 1
	it.decodeA.Decode(d, enableDecodeA)

	enableDecodeB := util.GetBit(d, 5) == 1
	it.decodeB.Decode(d>>5, enableDecodeB)

	Broker.Publish("PRGC_E", util.GetBit(d, 10))
	Broker.Publish("HLT", util.GetBit(d, 11))
	Broker.Publish("STPC_RST", util.GetBit(d, 12))
}
