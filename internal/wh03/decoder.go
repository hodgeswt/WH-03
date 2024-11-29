package wh03

import (
	"fmt"
)

type IDecoder interface {
	Decode(data int, enable bool)
}

type Decoder struct {
	Bitwidth int
	Outputs  map[int]string
}

func (it *Decoder) Decode(data int, enable bool) {
	if it.Outputs == nil {
		panic("Uninitialized decoder attempted to decode data")
	}

	for i := 0; i < (2 ^ it.Bitwidth); i++ {
		if it.Outputs[i] == "" {
			panic(fmt.Sprintf("Decoder found data for invalid output: %d", i))
		}

		if !enable {
			Broker.Publish(it.Outputs[i], 0)
		} else {
			v := (data >> i) & 1
			Broker.Publish(it.Outputs[i], v)

		}
	}
}
