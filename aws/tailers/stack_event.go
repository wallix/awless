package awstailers

import (
	"fmt"
	"io"
	"strings"
	"time"

	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
	"github.com/wallix/awless/aws/services"
)

type stackEventTailer struct {
	stackName        string
	follow           bool
	pollingFrequency time.Duration
	lastEventID      *string
	nbEvents         int
}

type stackEvents []*cloudformation.StackEvent

func NewCloudformationEventsTailer(stackName string, nbEvents int, enableFollow bool, frequency time.Duration) *stackEventTailer {
	return &stackEventTailer{stackName: stackName, follow: enableFollow, pollingFrequency: frequency, nbEvents: nbEvents}
}

func (t *stackEventTailer) Name() string {
	return "stack-events"
}

func (t *stackEventTailer) Tail(w io.Writer) error {
	cfn, ok := awsservices.CloudformationService.(*awsservices.Cloudformation)
	if !ok {
		return fmt.Errorf("invalid cloud service, expected awsservices.Cloudformation, got %T", awsservices.CloudformationService)
	}

	if t.pollingFrequency < 5*time.Second {
		return fmt.Errorf("invalid polling frequency: %s", t.pollingFrequency)
	}

	if !t.follow {
		if err := t.displayLastEvents(cfn, w); err != nil {
			return err
		}

		return nil
	}

	isDeploying, err := t.isStackBeingDeployed(cfn)
	if err != nil {
		return err
	}

	if !isDeploying {
		return fmt.Errorf("Stack %s not being deployed at the moment", t.stackName)
	}

	ticker := time.NewTicker(t.pollingFrequency)
	defer ticker.Stop()
	for range ticker.C {
		isDeploymentFinished, err := t.displayRelevantEvents(cfn, w)
		if err != nil || isDeploymentFinished {
			return err
		}
	}

	return nil
}

// get N latest events
func (t *stackEventTailer) getLatestEvents(cfn *awsservices.Cloudformation) (stackEvents, error) {
	params := &cloudformation.DescribeStackEventsInput{
		StackName: &t.stackName,
	}

	var stEvents stackEvents

	for {
		resp, err := cfn.DescribeStackEvents(params)
		if err != nil {
			return nil, err
		}

		for _, e := range resp.StackEvents {
			// if lastEventID == nil, then it's first run, and we just take first N events
			if t.lastEventID == nil && len(stEvents) >= t.nbEvents {
				return stEvents, nil
			}

			// if lastEventID found, then take all unseen events
			if t.lastEventID != nil && *e.EventId == *t.lastEventID {
				return stEvents, nil
			}
			stEvents = append(stEvents, e)
		}

		if resp.NextToken == nil {
			return stEvents, nil
		}

		params.NextToken = resp.NextToken
	}
}

func (t *stackEventTailer) displayLastEvents(cfn *awsservices.Cloudformation, w io.Writer) error {
	events, err := t.getLatestEvents(cfn)
	if err != nil {
		return err
	}

	if len(events) > 0 {
		t.lastEventID = events[0].EventId
		return events.printReverse(w)
	}

	return nil
}

func (t *stackEventTailer) isStackBeingDeployed(cfn *awsservices.Cloudformation) (bool, error) {
	stacks, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{StackName: &t.stackName})
	if err != nil {
		return false, err
	}

	if len(stacks.Stacks) == 0 {
		return false, fmt.Errorf("Stack not found")
	}

	return strings.HasSuffix(*stacks.Stacks[0].StackStatus, "_IN_PROGRESS"), nil
}

// get last N events relevant for current deployment in progress
func (t *stackEventTailer) getRelevantEvents(cfn *awsservices.Cloudformation) (stackEvents, bool, error) {
	params := &cloudformation.DescribeStackEventsInput{
		StackName: &t.stackName,
	}

	var stEvents stackEvents

	for {
		resp, err := cfn.DescribeStackEvents(params)
		if err != nil {
			return nil, false, err
		}

		for _, e := range resp.StackEvents {
			// if lastEventID == nil then it's first run of this method
			// if lastEventID == nil then it's not first run and print only new messages
			if t.lastEventID != nil && *e.EventId == *t.lastEventID {
				return stEvents, false, nil
			}
			stEvents = append(stEvents, e)
			// looking for the message which says that stack update or create started
			// making it as a first messages in the deployment events
			if *e.ResourceType == "AWS::CloudFormation::Stack" && strings.HasSuffix(*e.ResourceStatus, "_IN_PROGRESS") {
				return stEvents, false, nil
			}

			// if we found next message, then stack create/update completed successfully
			if *e.ResourceType == "AWS::CloudFormation::Stack" && strings.HasSuffix(*e.ResourceStatus, "_COMPLETE") {
				return stEvents, true, nil
			}

			// if we found next message, then stack create/update failed
			if *e.ResourceType == "AWS::CloudFormation::Stack" && strings.HasSuffix(*e.ResourceStatus, "_FAILED") {
				return stEvents, true, fmt.Errorf("Deployment failed with error: %s", *e.ResourceStatusReason)
			}
		}

		if resp.NextToken == nil {
			return stEvents, false, nil
		}

		params.NextToken = resp.NextToken
	}

}

func (t *stackEventTailer) displayRelevantEvents(cfn *awsservices.Cloudformation, w io.Writer) (bool, error) {
	events, deploymentFinished, err := t.getRelevantEvents(cfn)
	if !deploymentFinished && err != nil {
		return deploymentFinished, err
	}

	if len(events) > 0 {
		t.lastEventID = events[0].EventId
	}

	events.printReverse(w)

	return deploymentFinished, err
}

func coloredResourceStatus(str string) string {
	switch {
	case strings.HasSuffix(str, "_IN_PROGRESS"):
		return color.New(color.FgYellow).SprintFunc()(str)
	case strings.HasSuffix(str, "_COMPLETE"):
		return color.New(color.FgGreen).SprintFunc()(str)
	case strings.HasSuffix(str, "_FAILED"):
		return color.New(color.FgRed).SprintFunc()(str)
	default:
		return str
	}

}

// TODO: add filters for printed fields
func (e stackEvents) printReverse(w io.Writer) error {
	tab := tabwriter.NewWriter(w, 25, 8, 0, '\t', 0)
	for i := len(e) - 1; i >= 0; i-- {
		_, err := fmt.Fprintf(tab, "%s\t%s\t%s\t%s\t", e[i].Timestamp.Format(time.RFC3339), *e[i].LogicalResourceId, *e[i].ResourceType, coloredResourceStatus(*e[i].ResourceStatus))
		if err != nil {
			return err
		}
		fmt.Fprintln(tab)
	}
	tab.Flush()
	return nil
}
