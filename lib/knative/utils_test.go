package knative

import (
	"io/ioutil"
	"testing"

	"k8s.io/apimachinery/pkg/types"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	"github.com/ghodss/yaml"
	"github.com/iter8-tools/handler/base/experiment"
	"github.com/iter8-tools/handler/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetNamespacedNameForKsvc(t *testing.T) {
	var err error
	var exp *experiment.Experiment
	var nn *types.NamespacedName
	exp, err = (&experiment.Builder{}).FromFile(utils.CompletePath("../../", "testdata/experiment6.yaml")).Build()
	assert.NoError(t, err)

	nn, err = GetNamespacedNameForKsvc(exp)
	assert.Equal(t, *nn, types.NamespacedName{
		Namespace: "default",
		Name:      "sklearn-iris",
	})
	assert.NoError(t, err)

	exp, err = (&experiment.Builder{}).FromFile(utils.CompletePath("../../", "testdata/experiment2.yaml")).Build()
	assert.NoError(t, err)

	nn, err = GetNamespacedNameForKsvc(exp)
	assert.Nil(t, nn)
	assert.Error(t, err)
}

func TestRevisionPresentInKsvc(t *testing.T) {
	// get ksvc
	ksvcBytes, err := ioutil.ReadFile(utils.CompletePath("../../", "testdata/knative/onerevision.yaml"))
	ksvc := &servingv1.Service{}
	err = yaml.Unmarshal(ksvcBytes, ksvc)
	assert.NoError(t, err)

	// assert revision in ksvc
	assert.True(t, revisionPresentInKsvc("sample-application-v1", ksvc))
}
