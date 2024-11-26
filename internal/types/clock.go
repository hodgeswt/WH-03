package types

import "context"

type Clock struct {
    Freq int
}

func (it *Clock) Run(ctx context.Context) {

}

func (it *Clock) Stop() {

}
