package notification

import (
	"errors"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/handler/base"
	"github.com/iter8-tools/handler/utils"
	"github.com/sirupsen/logrus"
)

const (
	LIBRARY string = "notification"
)

var log *logrus.Logger

func init() {
	log = utils.GetLogger()
}

// MakeTask constructs a Task from a TaskMeta or returns an error if any.
func MakeTask(t *v2alpha2.TaskSpec) (base.Task, error) {
	switch t.Task {
	case LIBRARY + "/" + SLACK_TASK:
		return MakeSlackTask(t)
	// add additional tasks here options here
	default:
		return nil, errors.New("Unknown task: " + t.Task)
	}
}
