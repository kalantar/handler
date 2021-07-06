package tasks_test

import (
	"testing"

	"github.com/iter8-tools/handler/tasks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMakeTask(t *testing.T) {
	log := tasks.GetLogger()
	log.Info("hello world")
	assert.NotEmpty(t, log)
	tasks.SetLogLevel(logrus.InfoLevel)
}

func TestGetJsonBytes(t *testing.T) {
	// valid
	_, err := tasks.GetJSONBytes("https://httpbin.org/stream/1")
	assert.NoError(t, err)

	// invalid
	_, err = tasks.GetJSONBytes("https://httpbin.org/undef")
	assert.Error(t, err)
}

func TestPointers(t *testing.T) {
	assert.Equal(t, int32(1), *tasks.Int32Pointer(1))
	assert.Equal(t, float32(0.1), *tasks.Float32Pointer(0.1))
	assert.Equal(t, float64(0.1), *tasks.Float64Pointer(0.1))
	assert.Equal(t, "hello", *tasks.StringPointer("hello"))
	assert.Equal(t, false, *tasks.BoolPointer(false))
	assert.Equal(t, tasks.GET, *tasks.HTTPMethodPointer(tasks.GET))
}
