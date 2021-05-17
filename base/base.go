package base

import (
	"bytes"
	"context"
	"errors"
	"html/template"

	"github.com/iter8-tools/handler/utils"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = utils.GetLogger()
}

// Task defines common method signatures for every task.
type Task interface {
	Run(ctx context.Context) error
}

// Action is a slice of Tasks.
type Action []Task

// TaskMeta is common to all Tasks
type TaskMeta struct {
	Library string `json:"library" yaml:"library"`
	Task    string `json:"task" yaml:"task"`
}

// Run the given action.
func (a *Action) Run(ctx context.Context) error {
	for i := 0; i < len(*a); i++ {
		log.Info("------")
		err := (*a)[i].Run(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Tags supports string extrapolation using tags.
type Tags struct {
	M map[string]string
}

// Interpolate str using tags.
func (tags *Tags) Interpolate(str *string) (string, error) {
	if tags == nil || tags.M == nil { // return a copy of the string
		return *str, nil
	}
	var err error
	var templ *template.Template
	if templ, err = template.New("").Parse(*str); err == nil {
		buf := bytes.Buffer{}
		if err = templ.Execute(&buf, tags.M); err == nil {
			return string(buf.Bytes()), nil
		}
		log.Error("template execution error: ", err)
		return "", errors.New("cannot interpolate string")
	}
	log.Error("template creation error: ", err)
	return "", errors.New("cannot interpolate string")
}

// ContextKey is the type of key that will be used to index into context.
type ContextKey string
