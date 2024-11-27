//go:build unit
// +build unit

package types

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStepCounterBuffer_NoBuffer(t *testing.T) {
    testStepCounter := &StepCounter{}

    testStepCounter.Buffer("testKey", "testValue")

    assert.Equal(t, testStepCounter.buffer["testKey"], "testValue")
}

func TestStepCounterBuffer_BufferPopulated(t *testing.T) {
    testStepCounter := &StepCounter{buffer: map[string]string{"testKey": "oldValue"}}

    testStepCounter.Buffer("testKey", "testValue")

    assert.Equal(t, testStepCounter.buffer["testKey"], "testValue")
}

func TestStepCounterBuffer_BufferEmpty(t *testing.T) {
    testStepCounter := &StepCounter{buffer: map[string]string{}}

    testStepCounter.Buffer("testKey", "testValue")

    assert.Equal(t, testStepCounter.buffer["testKey"], "testValue")
}

func TestStepCounterUpdateState_Initial0(t *testing.T) {
    testStepCounter := &StepCounter{limit: 8}

    testStepCounter.UpdateState()

    assert.Equal(t, testStepCounter.state, 1)
}

func TestStepCounterUpdateState_InitialAtLimit(t *testing.T) {
    testStepCounter := &StepCounter{limit: 7, state: 7}

    testStepCounter.UpdateState()

    assert.Equal(t, testStepCounter.state, 0)
}

func TestStepCounterResetState(t *testing.T) {
    testStepCounter := &StepCounter{state: 7}

    testStepCounter.Reset()

    assert.Equal(t, testStepCounter.state, 0)
}

