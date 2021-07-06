package notification

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"strings"

// 	"github.com/iter8-tools/etc3/api/v2alpha2"
// 	"github.com/iter8-tools/handler/base"
// 	"github.com/iter8-tools/handler/experiment"
// 	"github.com/iter8-tools/handler/lib/common"
// 	"github.com/slack-go/slack"
// 	"github.com/spf13/viper"
// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/apimachinery/pkg/types"
// )

// const (
// 	// Slack2TaskName is the name of the task this file implements
// 	Slack2TaskName string = "slack2"
// )

// // Slack2TaskInputs is the object corresponding to the expcted inputs to the task
// type Slack2TaskInputs struct {
// 	Channel string `json:"channel" yaml:"channel"`
// 	Secret  string `json:"secret" yaml:"secret"`
// }

// // SlackTask encapsulates a command that can be executed.
// type Slack2Task struct {
// 	base.TaskMeta `json:",inline" yaml:",inline"`
// 	// If there are any additional inputs
// 	With SlackTaskInputs `json:"with" yaml:"with"`
// }

// // MakeSlack2Task converts an sampletask spec into an base.Task.
// func MakeSlack2Task(t *v2alpha2.TaskSpec) (base.Task, error) {
// 	if t.Task != LibraryName+"/"+SlackTaskName {
// 		return nil, fmt.Errorf("library and task need to be '%s' and '%s'", LibraryName, SlackTaskName)
// 	}
// 	var jsonBytes []byte
// 	var task base.Task
// 	// convert t to jsonBytes
// 	jsonBytes, err := json.Marshal(t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// convert jsonString to SlackTask
// 	task = &SlackTask{}
// 	err = json.Unmarshal(jsonBytes, &task)
// 	return task, err
// }

// // Run the task. This suppresses all errors so that the task will always succeed.
// // In this way, any failure does not cause failure of the enclosing experiment.
// func (t *Slack2Task) Run(ctx context.Context) error {
// 	t.internalRun(ctx)
// 	return nil
// }

// func (t *Slack2Task) ToHttpRequestTask() *common.HttpRequestTask {
// 	// curl -X POST -H 'Content-type: application/json' --data '{"text":"Hello, World!"}' https://hooks.slack.com/services/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
// 	body := `
// {
// 	text: "Hello, world."
// }
// `
// 	method := v2alpha2.POSTMethodType
// 	tSpec := &common.HttpRequestTask{
// 		TaskMeta: base.TaskMeta{
// 			Library: common.LibraryName,
// 			Task:    common.HttpRequestTaskName,
// 		},
// 		With: common.HttpRequestInputs{
// 			URL:    "http://hooks.slack.com/services/TQ3FN6N01/CU5FNKWCB/xoxb-819532226001-2019118254530-kXaCV8F0z6LK78tk0h4kgPol",
// 			Method: &method,
// 			Headers: []v2alpha2.NamedValue{{
// 				Name:  "Content-type",
// 				Value: "application/json",
// 			}},
// 			Body: &body,
// 		},
// 	}

// 	return tSpec
// }

// // Actual task runner
// func (t *Slack2Task) internalRun(ctx context.Context) error {
// 	// Called to execute the Task
// 	// Retrieve the experiment object (if needed)
// 	exp, err := experiment.GetExperimentFromContext(ctx)
// 	// exit with error if unable to retrieve experiment
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	log.Trace("experiment", exp)
// 	return t.postNotification(exp)
// }

// func (t *Slack2Task) postNotification(e *experiment.Experiment) error {
// 	token := t.getToken()
// 	if token == nil {
// 		return errors.New("Unable to find token")
// 	}
// 	log.Trace("token", t.getToken())
// 	api := slack.New(*token)
// 	channelID, timestamp, err := api.PostMessage(
// 		t.With.Channel,
// 		slack.MsgOptionBlocks(slack.NewSectionBlock(&slack.TextBlockObject{
// 			Type: slack.MarkdownType,
// 			// Text: Bold(Name(e)),
// 			Text: Bold(string(e.Spec.Strategy.TestingPattern) + " experiment on " + e.Spec.Target),
// 		}, nil, nil)),
// 		slack.MsgOptionAttachments(slack.Attachment{
// 			Blocks: slack.Blocks{
// 				BlockSet: []slack.Block{
// 					slack.NewSectionBlock(&slack.TextBlockObject{
// 						Type: slack.MarkdownType,
// 						Text: SlackMessage(e),
// 					}, nil, nil),
// 				},
// 			},
// 		}),
// 		slack.MsgOptionIconURL("https://avatars.githubusercontent.com/u/53243580?s=200&v=4"),
// 	)

// 	log.Trace("channelID", channelID)
// 	log.Trace("timestamp", timestamp)
// 	return err
// }

// // SlackMessage constructs the slack message to post
// func SlackMessage(e *experiment.Experiment) string {
// 	msg := []string{
// 		// Bold("Type: ") + Italic(string(e.Spec.Strategy.TestingPattern)),
// 		// Bold("Target: ") + Italic(e.Spec.Target),
// 		Bold("Name:") + Space + Italic(Name(e)),
// 		Bold("Versions:") + Space + Italic(Versions(e)),
// 		Bold("Stage:") + Space + Italic(Stage(e)),
// 		Bold("Winner:") + Space + Italic(Winner(e)),
// 	}

// 	if Failed(e) {
// 		msg = append(msg, Bold("Failed:")+Space+Italic("true"))
// 	}

// 	return strings.Join(msg, NewLine)
// }

// func (t *Slack2Task) getToken() *string {
// 	// get secret namespace and name
// 	namespace := viper.GetViper().GetString("experiment_namespace")
// 	var name string
// 	secretNN := t.With.Secret
// 	nn := strings.Split(secretNN, "/")
// 	if len(nn) == 1 {
// 		name = nn[0]
// 	} else {
// 		namespace = nn[0]
// 		name = nn[1]
// 	}
// 	log.Trace("namespace", namespace)
// 	log.Trace("name", name)

// 	secret := corev1.Secret{}
// 	err := experiment.GetTypedObject(&types.NamespacedName{Namespace: namespace, Name: name}, &secret)

// 	if err != nil {
// 		log.Error(err)
// 		return nil
// 	}
// 	token := string(secret.Data["token"])
// 	return &token
// }
