package wh03

import (
	"context"
	"fmt"

	"github.com/hodgeswt/WH-03/internal/util"
	"github.com/hodgeswt/utilw/pkg/logw"
)

type IAlu interface {
	ICircuit
}

type Alu struct {
	buffer map[string]int
}

func (it *Alu) Buffer(key string, data int) {
	logw.Debug("^Alu.Buffer")
	defer logw.Debug("$Alu.Buffer")

	if it.buffer == nil {
		it.buffer = map[string]int{}
	}

	it.buffer[key] = data
}

func (it *Alu) Reset() {
	logw.Debug("^$Alu.Reset")
}

func (it *Alu) Run(ctx context.Context) {
	logw.Debug("^Alu.Run")
	defer logw.Debug("$Alu.Run")

	da := Broker.Subscribe("A_D")
	db := Broker.Subscribe("B_D")
	op := Broker.Subscribe("Alu_OP")
	rst := Broker.Subscribe("RST")

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-da:
			logw.Infof("Alu received A reg data: %08b", dat)
			it.Buffer("A_D", dat)
			it.UpdateState()
		case dat := <-db:
			logw.Infof("Alu received B reg data: %08b", dat)
			it.Buffer("B_D", dat)
			it.UpdateState()
		case dat := <-op:
			logw.Infof("Alu received OP: %08b", dat)
			it.Buffer("B_D", dat)
			it.UpdateState()
		case dat := <-rst:
			logw.Infof("Alu received RST update %08b", dat)
			it.Reset()
		}
	}
}

func (it *Alu) UpdateState() {
	logw.Debug("^Alu.UpdateState")
	defer logw.Debug("$Alu.UpdateState")

	da := it.buffer["A_D"]
	db := it.buffer["B_D"]
	op := it.buffer["OP"]

	add := 0
	sub := 1
	lshift := 2
	rshift := 3

	dat := -1
	carry := 0

	switch op {
	case add:
		dat = da + db
		if dat > util.MaxBit {
			dat = util.MaxBit
			carry = 1
        }
        break
	case sub:
		dat = da - db
		if dat < 0 {
			dat = 0
			carry = 1
		}
		break
	case lshift:
		dat = da << db
		break
	case rshift:
		dat = da >> db
		break
	}

	if dat == -1 {
		panic(fmt.Sprintf("Unexpected ALU result from data: %08b, %08b, op: %08b", da, db, op))
	}

	if it.buffer["Accumulator_OE"] == 1 {
		Broker.Publish("D", dat)
	}

	Broker.Publish("Alu_D", dat)

	gt := 0
	lt := 0
	eq := 0

	if da > db {
		gt = 1
	}

	if da < db {
		lt = 1
	}

	if da == db {
		eq = 1
	}

	Broker.Publish("FD", gt<<0|lt<<1|eq<<2|carry<<3)
}
