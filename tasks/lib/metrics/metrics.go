package metrics

import (
	"errors"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/handler/tasks"
	"github.com/sirupsen/logrus"
)

const (

	// LibraryName is the name of the task library this package implements
	LibraryName string = "metrics"
)

// Declare logger only once per package (in any file belonging to that package)
var log *logrus.Logger

func init() {
	// always use logger from utils
	// init logger once per package (in any file belonging to that package)
	log = tasks.GetLogger()
}

// MakeTask constructs a Task from a TaskMeta or returns an error if any.
func MakeTask(t *v2alpha2.TaskSpec) (tasks.Task, error) {
	switch t.Task {
	case LibraryName + "/" + CollectTaskName:
		bt, err := MakeCollect(t)
		return bt, err
	default:
		return nil, errors.New("Unknown task: " + t.Task)
	}
}
