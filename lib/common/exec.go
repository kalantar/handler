package common

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/handler/base"
	"github.com/iter8-tools/handler/base/experiment"
)

const (
	// ExecTaskName is the name of this file implements
	ExecTaskName string = "exec"
)

// ExecInputs contain the name and arguments of the command to be executed.
type ExecInputs struct {
	Cmd                  string        `json:"cmd" yaml:"cmd"`
	Args                 []interface{} `json:"args,omitempty" yaml:"args,omitempty"`
	DisableInterpolation bool          `json:"disableInterpolation,omitempty" yaml:"disableInterpolation,omitempty"`
}

// ExecTask encapsulates a command that can be executed.
type ExecTask struct {
	base.TaskMeta `json:",inline" yaml:",inline"`
	With          ExecInputs `json:"with" yaml:"with"`
}

// Run the command.
func (t *ExecTask) Run(ctx context.Context) error {
	exp, err := experiment.GetExperimentFromContext(ctx)
	if err == nil {
		inputArgs := make([]string, len(t.With.Args))
		for i := 0; i < len(inputArgs); i++ {
			inputArgs[i] = fmt.Sprint(t.With.Args[i])
		}
		log.Trace(inputArgs)
		var args []string
		if t.With.DisableInterpolation {
			args = inputArgs
		} else {
			args, err = exp.Interpolate(inputArgs)
		}
		if err == nil {
			log.Trace("interpolated args: ", args)
			cmd := exec.Command(t.With.Cmd, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			log.Info("Running task: " + cmd.String())
			log.Trace(args)
			err = cmd.Run()
		}
	}
	if err != nil {
		log.Error(err)
	}
	return err
}

// MakeExec converts an exec task spec into an exec task.
func MakeExec(t *v2alpha2.TaskSpec) (base.Task, error) {
	if t.Task != LibraryName+"/"+ExecTaskName {
		return nil, fmt.Errorf("library and task need to be '%s' and '%s'", LibraryName, ExecTaskName)
	}
	var err error
	var jsonBytes []byte
	var et base.Task
	// convert t to jsonBytes
	jsonBytes, err = json.Marshal(t)
	// convert jsonString to ExecTask
	if err == nil {
		et = &ExecTask{}
		err = json.Unmarshal(jsonBytes, &et)
	}
	return et, err
}
