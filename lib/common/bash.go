package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/handler/base"
)

const (
	// BashTaskName is the name of the bash task
	BashTaskName string = "bash"
)

// BashInputs contain the name and arguments of the command to be executed.
type BashInputs struct {
	Script string `json:"script" yaml:"script"`
}

// BashTask encapsulates a command that can be executed.
type BashTask struct {
	base.TaskMeta `json:",inline" yaml:",inline"`
	With          BashInputs `json:"with" yaml:"with"`
}

// MakeBashTask converts an exec task spec into an exec task.
func MakeBashTask(t *v2alpha2.TaskSpec) (base.Task, error) {
	if t.Task != LibraryName+"/"+BashTaskName {
		return nil, fmt.Errorf("library and task need to be '%s' and '%s'", LibraryName, BashTaskName)
	}
	var jsonBytes []byte
	var task base.Task
	// convert t to jsonBytes
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	// convert jsonString to ExecTask
	task = &BashTask{}
	err = json.Unmarshal(jsonBytes, &task)
	return task, err
}

// Run the command.
func (t *BashTask) Run(ctx context.Context) error {
	tags := base.GetDefaultTags(ctx)
	log.Tracef("tags: %v", *tags)

	// interpolate - replaces placeholders in the script with values
	script, err := tags.Interpolate(&t.With.Script)

	log.Trace(script)
	args := []string{"-c", script}

	cmd := exec.Command("/bin/bash", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Info("Running task: " + cmd.String())
	log.Trace(args)
	err = cmd.Run()

	return err
}
