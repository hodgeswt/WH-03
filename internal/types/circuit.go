package types

import "context"

type ICircuit interface {
    Run(ctx context.Context)
}
