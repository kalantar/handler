package notification

import (
	"errors"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/handler/tasks"
	"github.com/sirupsen/logrus"
)

const (

	// LibraryName is the name of the task library this package implements
	LibraryName string = "notification"
)

var log *logrus.Logger

func init() {
	log = tasks.GetLogger()
}

// MakeTask constructs a Task from a TaskMeta or returns an error if any.
func MakeTask(t *v2alpha2.TaskSpec) (tasks.Task, error) {
	switch t.Task {
	case LibraryName + "/" + SlackTaskName:
		return MakeSlackTask(t)
	// add additional tasks here
	default:
		return nil, errors.New("Unknown task: " + t.Task)
	}
}
