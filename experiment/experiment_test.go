package experiment

import (
	"context"
	"testing"

	"github.com/iter8-tools/etc3/api/v2alpha1"
	"github.com/iter8-tools/handler/base"
	"github.com/iter8-tools/handler/utils"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuildErrorGarbageYAML(t *testing.T) {
	_, err := (&Builder{}).FromFile(utils.CompletePath("../", "testdata/garbage.yaml")).Build()
	assert.Error(t, err)
}

func TestInvalidAction(t *testing.T) {
	_, err := (&Builder{}).FromFile(utils.CompletePath("../", "testdata/experiment3.yaml")).Build()
	assert.Error(t, err)
}

func TestInvalidActions(t *testing.T) {
	_, err := (&Builder{}).FromFile(utils.CompletePath("../", "testdata/experiment5.yaml")).Build()
	assert.Error(t, err)
}

func TestStringAction(t *testing.T) {
	_, err := (&Builder{}).FromFile(utils.CompletePath("../", "testdata/experiment9.yaml")).Build()
	assert.Error(t, err)
}

func TestGetRecommendedBaseline(t *testing.T) {
	var err error
	var exp *Experiment
	var b string
	exp, err = (&Builder{}).FromFile(utils.CompletePath("../", "testdata/experiment6.yaml")).Build()
	assert.NoError(t, err)
	b, err = exp.GetRecommendedBaseline()
	assert.NoError(t, err)
	assert.Equal(t, "default", b)

	exp, err = (&Builder{}).FromFile(utils.CompletePath("../", "testdata/experiment2.yaml")).Build()
	assert.NoError(t, err)
	b, err = exp.GetRecommendedBaseline()
	assert.Error(t, err)

	exp = nil
	b, err = exp.GetRecommendedBaseline()
	assert.Error(t, err)
}

func TestGetExperimentFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), base.ContextKey("experiment"), "hello world")
	_, err := GetExperimentFromContext(ctx)
	assert.Error(t, err)

	_, err = GetExperimentFromContext(context.Background())
	assert.Error(t, err)

	ctx = context.WithValue(context.Background(), base.ContextKey("experiment"), &Experiment{
		Experiment: v2alpha1.Experiment{
			TypeMeta:   v1.TypeMeta{},
			ObjectMeta: v1.ObjectMeta{},
			Spec:       v2alpha1.ExperimentSpec{},
			Status:     v2alpha1.ExperimentStatus{},
		},
	})

	exp, err := GetExperimentFromContext(ctx)
	assert.NotNil(t, exp)
	assert.NoError(t, err)
}

func TestExtrapolate(t *testing.T) {
	var err error
	var exp *Experiment
	var b string
	exp, err = (&Builder{}).FromFile(utils.CompletePath("../", "testdata/experiment6.yaml")).Build()
	assert.NoError(t, err)
	b, err = exp.GetRecommendedBaseline()
	assert.NoError(t, err)
	assert.Equal(t, "default", b)

	args, err := exp.Extrapolate(nil)
	assert.Empty(t, args)
	assert.NoError(t, err)

	args, err = exp.Extrapolate([]string{"hello-world", "hello {{ .revision }} world", "hello {{ .omg }} world"})
	assert.Equal(t, []string{"hello-world", "hello revision1 world", "hello  world"}, args)
	assert.NoError(t, err)

	args, err = exp.Extrapolate([]string{"hello-world", "hello {{ .revision }} world", "hello {{ range .ForEver .omg }} world"})
	assert.Error(t, err)
}

// func TestExtrapolateWithoutHandlerStanza(t *testing.T) {
// 	var e *Experiment = &Experiment{}
// 	assert.NoError(t, e.extrapolate())

// 	_, err := e.GetAction("start")
// 	assert.Error(t, err)

// 	e.Spec.Strategy.Handlers = &Handlers{}
// 	_, err = e.GetAction("start")
// 	assert.Error(t, err)
// }

// func TestExtrapolateWithoutRecommendedBaseline(t *testing.T) {
// 	var e *Experiment = &Experiment{
// 		Experiment: *iter8.NewExperiment("default", "default").Build(),
// 	}
// 	e.Spec.Strategy.Handlers = &Handlers{
// 		Actions: &ActionMap{},
// 	}
// 	assert.Error(t, e.extrapolate())
// }

// func TestExtrapolateWithoutVersionTags(t *testing.T) {
// 	var iter8Experiment = *iter8.NewExperiment("some", "exp").Build()
// 	spec := Spec{}
// 	spec.VersionInfo = &iter8.VersionInfo{
// 		Baseline: iter8.VersionDetail{
// 			Name:         "default",
// 			Tags:         nil,
// 			WeightObjRef: &v1.ObjectReference{},
// 		},
// 		Candidates: []iter8.VersionDetail{
// 			{
// 				Name:         "canary",
// 				Tags:         nil,
// 				WeightObjRef: &v1.ObjectReference{},
// 			},
// 		},
// 	}
// 	var e *Experiment = &Experiment{
// 		Experiment: iter8Experiment,
// 		Spec:       spec,
// 	}
// 	x := "default"
// 	e.Status.RecommendedBaseline = &x

// 	e.Spec.Strategy.Handlers = &Handlers{
// 		Actions: &ActionMap{},
// 	}
// 	assert.NoError(t, e.extrapolate())
// }

// func TestGetVersionInfo(t *testing.T) {
// 	var iter8Experiment = *iter8.NewExperiment("some", "exp").Build()
// 	spec := Spec{}
// 	spec.VersionInfo = &iter8.VersionInfo{
// 		Baseline: iter8.VersionDetail{
// 			Name:         "default",
// 			Tags:         nil,
// 			WeightObjRef: &v1.ObjectReference{},
// 		},
// 		Candidates: []iter8.VersionDetail{
// 			{
// 				Name:         "canary",
// 				Tags:         nil,
// 				WeightObjRef: &v1.ObjectReference{},
// 			},
// 		},
// 	}
// 	var e *Experiment = &Experiment{
// 		Experiment: iter8Experiment,
// 		Spec:       spec,
// 	}
// 	x := "default"
// 	e.Status.RecommendedBaseline = &x
// 	_, err := e.getVersionDetail("default")
// 	assert.NoError(t, err)

// 	e.Spec.Strategy.Handlers = &Handlers{
// 		Actions: &ActionMap{},
// 	}
// 	assert.NoError(t, e.extrapolate())

// 	_, err = e.getVersionDetail("canary")
// 	assert.NoError(t, err)

// 	_, err = e.getVersionDetail("random")
// 	assert.Error(t, err)
// }
