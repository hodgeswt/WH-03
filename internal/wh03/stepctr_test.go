//go:build unit
// +build unit

package wh03

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStepCounterBuffer_NoBuffer(t *testing.T) {
    testStepCounter := &StepCounter{}

    testStepCounter.Buffer("testKey", 1)

    assert.Equal(t, testStepCounter.buffer["testKey"], 1)
}

func TestStepCounterBuffer_BufferPopulated(t *testing.T) {
    testStepCounter := &StepCounter{buffer: map[string]int{"testKey": 0}}

    testStepCounter.Buffer("testKey", 1)

    assert.Equal(t, testStepCounter.buffer["testKey"], 1)
}

func TestStepCounterBuffer_BufferEmpty(t *testing.T) {
    testStepCounter := &StepCounter{buffer: map[string]int{}}

    testStepCounter.Buffer("testKey", 1)

    assert.Equal(t, testStepCounter.buffer["testKey"], 1)
}

func TestStepCounterUpdateState_Initial0(t *testing.T) {
    testStepCounter := &StepCounter{Limit: 8}

    testStepCounter.UpdateState()

    assert.Equal(t, testStepCounter.state, 1)
}

func TestStepCounterUpdateState_InitialAtLimit(t *testing.T) {
    testStepCounter := &StepCounter{Limit: 7, state: 7}

    testStepCounter.UpdateState()

    assert.Equal(t, testStepCounter.state, 0)
}

func TestStepCounterResetState(t *testing.T) {
    testStepCounter := &StepCounter{state: 7}

    testStepCounter.Reset()

    assert.Equal(t, testStepCounter.state, 0)
}

