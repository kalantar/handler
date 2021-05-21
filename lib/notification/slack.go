package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/iter8-tools/etc3/api/v2alpha2"
	"github.com/iter8-tools/handler/base"
	"github.com/iter8-tools/handler/experiment"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	SLACK_TASK     string = "slack"
	SLACK_ENDPOINT string = "https://slack.com/api/chat.postMessage"
)

// SlackTaskInputs is the object corresponding to the expcted inputs to the task
type SlackTaskInputs struct {
	Channel string `json:"channel" yaml:"channel"`
	Secret  string `json:"secret" yaml:"secret"`
}

// SlackTask encapsulates a command that can be executed.
type SlackTask struct {
	base.TaskMeta `json:",inline" yaml:",inline"`
	// If there are any additional inputs
	With SlackTaskInputs `json:"with" yaml:"with"`
}

// MakeSlackTask converts an sampletask spec into an base.Task.
func MakeSlackTask(t *v2alpha2.TaskSpec) (base.Task, error) {
	if t.Task != LIBRARY+"/"+SLACK_TASK {
		return nil, errors.New(fmt.Sprintf("library and task need to be '%s' and '%s'", LIBRARY, SLACK_TASK))
	}
	var jsonBytes []byte
	var task base.Task
	// convert t to jsonBytes
	jsonBytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	// convert jsonString to SlackTask
	task = &SlackTask{}
	err = json.Unmarshal(jsonBytes, &task)
	return task, err
}

// Run the task. This suppresses all errors so that the task will always succeed.
// In this way, any failure does not cause failure of the enclosing experiment.
func (t *SlackTask) Run(ctx context.Context) error {
	t.internalRun(ctx)
	return nil
}

// Actual task runner
func (t *SlackTask) internalRun(ctx context.Context) error {
	// Called to execute the Task
	// Retrieve the experiment object (if needed)
	exp, err := experiment.GetExperimentFromContext(ctx)
	// exit with error if unable to retrieve experiment
	if err != nil {
		log.Error(err)
		return err
	}
	log.Trace("experiment", exp)
	return t.postNotification(exp)
}

func (t *SlackTask) postNotification(e *experiment.Experiment) error {
	token := t.getToken()
	if token == nil {
		return errors.New("Unable to find token")
	}
	log.Trace("token", t.getToken())
	api := slack.New(*token)
	channelID, timestamp, err := api.PostMessage(
		t.With.Channel,
		slack.MsgOptionBlocks(slack.NewSectionBlock(&slack.TextBlockObject{
			Type: slack.MarkdownType,
			Text: Bold(Name(e)),
		}, nil, nil)),
		slack.MsgOptionAttachments(slack.Attachment{
			Blocks: slack.Blocks{
				BlockSet: []slack.Block{
					slack.NewSectionBlock(&slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: SlackMessage(e),
					}, nil, nil),
				},
			},
		}),
		slack.MsgOptionIconURL("https://avatars.githubusercontent.com/u/53243580?s=200&v=4"),
	)

	log.Trace("channelID", channelID)
	log.Trace("timestamp", timestamp)
	return err
}

func SlackMessage(e *experiment.Experiment) string {
	msg := []string{
		"Type: " + Italic(string(e.Spec.Strategy.TestingPattern)),
		"Target: " + Italic(e.Spec.Target),
		"Versions: " + Italic(Versions(e)),
		"Stage: " + Italic(Stage(e)),
		"Winner: " + Italic(Winner(e)),
	}

	if Failed(e) {
		msg = append(msg, "Failed: "+Italic("true"))
	}

	return strings.Join(msg, NewLine())
}

func Name(e *experiment.Experiment) string {
	ns := e.Namespace
	if len(ns) == 0 {
		ns = "default"
	}
	return ns + "/" + e.Name
}

func Versions(e *experiment.Experiment) string {
	versions := make([]string, 0)
	if e.Spec.VersionInfo != nil {
		versions = append(versions, e.Spec.VersionInfo.Baseline.Name)
		for _, c := range e.Spec.VersionInfo.Candidates {
			versions = append(versions, c.Name)
		}
	}
	return strings.Join(versions, ", ")
}

func Stage(e *experiment.Experiment) string {
	stage := v2alpha2.ExperimentStageWaiting
	if e.Status.Stage != nil {
		stage = *e.Status.Stage
	}
	return string(stage)
}

func Winner(e *experiment.Experiment) string {
	winner := "not found"
	if e.Status.Analysis != nil &&
		e.Status.Analysis.WinnerAssessment != nil {
		if e.Status.Analysis.WinnerAssessment.Data.WinnerFound {
			winner = *e.Status.Analysis.WinnerAssessment.Data.Winner
		}
	}
	return winner
}

func Failed(e *experiment.Experiment) bool {
	// use !.. IsFalse() to allow undefined value => true
	return !e.Status.GetCondition(v2alpha2.ExperimentConditionExperimentFailed).IsFalse()
}

func Bold(text string) string {
	return "*" + text + "*"
}

func Italic(text string) string {
	return "_" + text + "_"
}

func NewLine() string {
	return "\n"
}

func (t *SlackTask) getToken() *string {
	// get secret namespace and name
	namespace := viper.GetViper().GetString("experiment_namespace")
	var name string
	secretNN := t.With.Secret
	nn := strings.Split(secretNN, "/")
	if len(nn) == 1 {
		name = nn[0]
	} else {
		namespace = nn[0]
		name = nn[1]
	}
	log.Trace("namespace", namespace)
	log.Trace("name", name)

	secret := corev1.Secret{}
	err := experiment.GetTypedObject(&types.NamespacedName{Namespace: namespace, Name: name}, &secret)

	if err != nil {
		log.Error(err)
		return nil
	}
	token := string(secret.Data["token"])
	return &token
}
