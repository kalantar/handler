package tasks_test

import (
	"context"
	"testing"

	"github.com/iter8-tools/handler/tasks"
	"github.com/stretchr/testify/assert"
)

func init() {
	log = tasks.GetLogger()
}

func TestWithExperiment(t *testing.T) {
	exp, err := (&tasks.Builder{}).FromFile(tasks.CompletePath("../", "testdata/experiment1.yaml")).Build()
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), tasks.ContextKey("experiment"), exp)
	tags := tasks.GetDefaultTags(ctx)

	testStr := []string{
		"{{.this.apiVersion}}",
		"{{.this.metadata.name}}",
		"{{.this.spec.duration.intervalSeconds}}",
		"{{(index .this.spec.versionInfo.baseline.variables 0).value}}",
		"{{.this.status.versionRecommendedForPromotion}}",
	}
	expectedOut := []string{
		"iter8.tools/v2alpha2",
		"sklearn-iris-experiment-1",
		"15",
		"revision1",
		"default",
	}

	for i, in := range testStr {
		out, err := tags.Interpolate(&in)
		assert.NoError(t, err)
		assert.Equal(t, expectedOut[i], out)
	}
}
