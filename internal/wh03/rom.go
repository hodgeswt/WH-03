package wh03

import (
	"context"
	"encoding/binary"
	"os"

	"github.com/hodgeswt/utilw/pkg/logw"
)

type IRom interface {
	ICircuit
	Load(path string) error
	Dump(path string) error
}

type Rom struct {
	Size   int
	buffer map[string]int
	state  []int64
	inst   int
	stpc   int
}

func (it *Rom) Buffer(key string, data int) {
	logw.Debug("^Rom.Buffer")
	defer logw.Debug("$Rom.Buffer")

	if it.buffer == nil {
		it.buffer = map[string]int{}
	}

	it.buffer[key] = data
}

func (it *Rom) Dump(path string) error {
	logw.Debug("^Rom.Dump")
	defer logw.Debug("$Rom.Dump")

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, it.state)

	return nil
}

func (it *Rom) Load(path string) error {
	logw.Debug("^Rom.Load")
	defer logw.Debug("$Rom.Load")

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	it.state = make([]int64, it.Size)

	err = binary.Read(file, binary.LittleEndian, &it.state)

	if err != nil {
		return err
	}

	return nil
}

func (it *Rom) Reset() {
	logw.Debug("^$Rom.Reset")

	it.buffer = map[string]int{}
}

func (it *Rom) Run(ctx context.Context) {
	logw.Debug("^Rom.Run")
	defer logw.Debug("$Rom.Run")

	stpc := Broker.Subscribe("STPC")
	inst := Broker.Subscribe("INST")
	rst := Broker.Subscribe("RST")

	for {
		select {
		case <-ctx.Done():
			return
		case dat := <-stpc:
			it.Buffer("STPC", dat)
			it.UpdateState()
		case dat := <-inst:
			it.Buffer("INST", dat)
			it.UpdateState()
		case dat := <-rst:
			logw.Infof("Rom received RST update %08b", dat)
			it.Reset()
		}
	}
}

func (it *Rom) UpdateState() {
	logw.Debug("^Rom.UpdateState")
	defer logw.Debug("$Rom.UpdateState")

	addr := (it.buffer["STPC"] << 8) | it.buffer["INST"]

	Broker.Publish("Rom_D", addr)
}
