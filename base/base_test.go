package base

import (
	"context"
	"testing"

	"github.com/iter8-tools/handler/experiment"
	"github.com/iter8-tools/handler/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	log = utils.GetLogger()
}

// func TestWithVersionRecommendedForPromotion(t *testing.T) {
// 	var data []byte
// 	data, err := ioutil.ReadFile(filepath.Join("..", "testdata", "experiment1.yaml"))
// 	assert.NoError(t, err)
// 	exp := &v2alpha2.Experiment{}
// 	err = yaml.Unmarshal(data, exp)
// 	assert.NoError(t, err)
// 	tags := interpolation.NewTags().WithRecommendedVersionForPromotion(exp)
// 	assert.Equal(t, "revision1", tags.M["revision"])
// }

// func TestWithOutVersionRecommendedForPromotion(t *testing.T) {
// 	var data []byte
// 	data, err := ioutil.ReadFile(filepath.Join("..", "testdata", "experiment1-norecommended.yaml"))
// 	assert.NoError(t, err)
// 	exp := &v2alpha2.Experiment{}
// 	err = yaml.Unmarshal(data, exp)
// 	assert.NoError(t, err)
// 	tags := interpolation.NewTags().WithRecommendedVersionForPromotion(exp)
// 	assert.NotContains(t, tags.M, "revision1")
// 	// assert.Equal(t, "revision1", tags.M["revision"])
// }

func TestWithExperiment(t *testing.T) {
	exp, err := (&experiment.Builder{}).FromFile(utils.CompletePath("../", "testdata/experiment1.yaml")).Build()
	assert.NoError(t, err)
	ctx := context.WithValue(context.Background(), utils.ContextKey("experiment"), exp)
	tags := GetDefaultTags(ctx)

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
