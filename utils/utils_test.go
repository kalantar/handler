package utils

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMakeTask(t *testing.T) {
	log := GetLogger()
	log.Info("hello world")
	assert.NotEmpty(t, log)
	SetLogLevel(logrus.InfoLevel)
}

func TestGetJsonBytes(t *testing.T) {
	// valid
	_, err := GetJSONBytes("https://httpbin.org/stream/1")
	assert.NoError(t, err)

	// invalid
	_, err = GetJSONBytes("https://httpbin.org/undef")
	assert.Error(t, err)
}

func TestPointers(t *testing.T) {
	assert.Equal(t, int32(1), *Int32Pointer(1))
	assert.Equal(t, float32(0.1), *Float32Pointer(0.1))
	assert.Equal(t, float64(0.1), *Float64Pointer(0.1))
	assert.Equal(t, "hello", *StringPointer("hello"))
	assert.Equal(t, false, *BoolPointer(false))
	assert.Equal(t, GET, *HTTPMethodPointer(GET))
}
