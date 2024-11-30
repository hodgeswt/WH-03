//go:build unit
// +build unit

package wh03

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPubSub_InitNoServer(t *testing.T) {
	pubsub := &PubSub{}

	pubsub.Init(10, false)

	assert.Equal(t, 10, pubsub.bufferSize)
	assert.Nil(t, pubsub.server)
}

func TestPubSub_InitServer(t *testing.T) {
	pubsub := &PubSub{}

	pubsub.Init(10, true)

	pubsub.Close()
	assert.Equal(t, 10, pubsub.bufferSize)
	assert.NotNil(t, pubsub.server)
}

func TestPubSub_Subscribe(t *testing.T) {
	pubsub := &PubSub{}
	pubsub.Init(10, false)

	ch := pubsub.Subscribe("test")

	assert.NotNil(t, ch)
}

func TestPubSub_Publish(t *testing.T) {
	pubsub := &PubSub{}
	pubsub.Init(10, false)

	ch := pubsub.Subscribe("test")
    done := make(chan struct{})
	dat := -1

	go func() {
		for {
			select {
			case d := <-ch:
				dat = d
                done <- struct{}{}
                close(done)
                return
			default:
				pubsub.Publish("test", 10)
			}

		}

	}()

    <-done

    select{
        case <-done:
            assert.Equal(t, 10, dat)
        case <-time.After(10 * time.Second):
            assert.FailNow(t, "nothing received from channel")
    }

}
