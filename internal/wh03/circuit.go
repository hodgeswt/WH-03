package wh03

import "context"

type ICircuit interface {
	Run(ctx context.Context)
	Buffer(key string, data int)
	UpdateState()
	Reset()
}
