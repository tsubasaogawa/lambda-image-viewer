package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcShutterSpeed(t *testing.T) {
	ss, _ := calcShutterSpeed(0)
	assert.Equal(t, "1/1", ss)
	ss, _ = calcShutterSpeed(1)
	assert.Equal(t, "1/2", ss)
	ss, _ = calcShutterSpeed(2)
	assert.Equal(t, "1/4", ss)
	_, err := calcShutterSpeed(-1)
	assert.Error(t, err)
	ss, _ = calcShutterSpeed(0.5)
	assert.Equal(t, "1/1", ss)
	ss, _ = calcShutterSpeed(10)
	assert.Equal(t, "1/1024", ss)
}
