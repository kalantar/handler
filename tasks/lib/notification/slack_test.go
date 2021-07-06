package notification

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/handler/tasks"
	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestMakeTask(t *testing.T) {
	channel, _ := json.Marshal("channel")
	secret, _ := json.Marshal("default/slack-secret")
	task, err := MakeTask(&v2alpha2.TaskSpec{
		Task: LibraryName + "/" + SlackTaskName,
		With: map[string]apiextensionsv1.JSON{
			"channel": {Raw: channel},
			"secret":  {Raw: secret},
		},
	})
	assert.NotEmpty(t, task)
	assert.NoError(t, err)
	assert.Equal(t, "channel", task.(*SlackTask).With.Channel)
	assert.Equal(t, "default/slack-secret", task.(*SlackTask).With.Secret)
}

type test struct {
	fileName         string
	expectedName     string
	expectedVersions string
	expectedStage    string
	expectedWinner   string
	expectedFailure  bool
}

const (
	winnerNotFound string = "not found"
)

// table driven tests
var tests = []test{
	// Conformance Test (1 versions), success, winner
	{fileName: "slack1.yaml", expectedName: "default/conformance-exp", expectedVersions: "productpage-v1", expectedStage: "Completed", expectedWinner: "productpage-v1", expectedFailure: false},
	// A/B test  (2 versions), failed
	{fileName: "slack2.yaml", expectedName: "default/quickstart-exp", expectedVersions: "productpage-v1, productpage-v2", expectedStage: "Completed", expectedWinner: winnerNotFound, expectedFailure: false},
	// A/B/n Test (3 versions), --> no analysis (winner); no failure, no stage
	{fileName: "slack3.yaml", expectedName: "default/abn-exp", expectedVersions: "productpage-v1, productpage-v2, productpage-v3", expectedStage: "Waiting", expectedWinner: winnerNotFound, expectedFailure: false},
}

func TestExperiment(t *testing.T) {
	for _, tc := range tests {

		exp, err := (&tasks.Builder{}).FromFile(filepath.Join("..", "..", "..", "testdata", "notification", tc.fileName)).Build()
		assert.NoError(t, err)
		msg := SlackMessage(exp)
		assert.Contains(t, Name(exp), tc.expectedName)
		assert.Contains(t, msg, "*Versions:* _"+tc.expectedVersions+"_\n")
		assert.Contains(t, msg, "*Winner:* _"+tc.expectedWinner+"_")
		if tc.expectedFailure {
			assert.Contains(t, msg, "Failed")
		}
	}
}
