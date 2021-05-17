package notification

import (
	"encoding/json"
	"testing"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestMakeTask(t *testing.T) {
	channel, _ := json.Marshal("channel")
	secret, _ := json.Marshal("default/slack-secret")
	task, err := MakeTask(&v2alpha2.TaskSpec{
		Task: LIBRARY + "/" + SLACK_TASK,
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
