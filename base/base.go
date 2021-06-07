package base

import (
	"bytes"
	"context"
	"errors"
	"html/template"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/handler/utils"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
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
	M map[string]interface{}
}

// NewTags creates an empty instance of Tags
func NewTags() Tags {
	return Tags{M: make(map[string]interface{})}
}

// WithSecret adds the fields in secret to tags
func (tags Tags) WithSecret(secret *corev1.Secret) Tags {
	if secret != nil {
		for n, v := range secret.Data {
			tags.M[n] = string(v)
		}
	}
	return tags
}

// With adds obj to tags
func (tags Tags) With(label string, obj interface{}) Tags {
	if obj != nil {
		tags.M[label] = obj
	}
	return tags
}

// WithRecommendedVersionForPromotion adds variables from versionDetail of version recommended for promotion
func (tags Tags) WithRecommendedVersionForPromotion(exp *v2alpha2.Experiment) Tags {
	if exp == nil || exp.Status.VersionRecommendedForPromotion == nil {
		log.Warn("no version recommended for promotion")
		return tags
	}

	versionRecommendedForPromotion := *exp.Status.VersionRecommendedForPromotion
	if exp.Spec.VersionInfo == nil {
		log.Warnf("No version details found for version recommended for promotion: %s", versionRecommendedForPromotion)
		return tags
	}

	var versionDetail *v2alpha2.VersionDetail = nil
	if exp.Spec.VersionInfo.Baseline.Name == versionRecommendedForPromotion {
		versionDetail = &exp.Spec.VersionInfo.Baseline
	} else {
		for _, v := range exp.Spec.VersionInfo.Candidates {
			if v.Name == versionRecommendedForPromotion {
				versionDetail = &v
				break
			}
		}
	}
	if versionDetail == nil {
		log.Warnf("No version details found for version recommended for promotion: %s", versionRecommendedForPromotion)
		return tags
	}

	// get the variable values from the (recommended) versionDetail
	tags.M["name"] = versionDetail.Name
	for _, v := range versionDetail.Variables {
		tags.M[v.Name] = v.Value
	}

	return tags
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
