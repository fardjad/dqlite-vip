package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitWithTimeout(t *testing.T) {
	var wg WaitGroupWithTimeout

	println(time.Now().String())

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Millisecond)
	}()
	assert.True(t, wg.WaitWithTimeout(1*time.Second))

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
	}()
	assert.False(t, wg.WaitWithTimeout(1*time.Millisecond))
}
