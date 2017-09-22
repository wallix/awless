package awstailers

import (
	"fmt"
	"io"
	"strings"
	"time"

	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/wallix/awless/aws/services"
)

type stackEventTailer struct {
	stackName        string
	refresh          bool
	watchDeployment  bool
	refreshFrequency time.Duration
	lastEventID      *string
	nbEvents         int
}

type stackEvents []*cloudformation.StackEvent

func (e stackEvents) print(w io.Writer) error {
	tab := tabwriter.NewWriter(w, 35, 8, 0, '\t', 0)
	for i, event := range e {
		_, err := fmt.Fprintf(tab, "%s\t%s\t%s\t%s\t%d", event.Timestamp.Format(time.RFC3339), *event.LogicalResourceId, *event.ResourceType, *event.ResourceStatus, i)
		if err != nil {
			return err
		}
		fmt.Fprintln(tab)
	}
	tab.Flush()
	return nil
}

func (e stackEvents) printReverse(w io.Writer) error {
	tab := tabwriter.NewWriter(w, 35, 8, 0, '\t', 0)
	for i := len(e) - 1; i >= 0; i-- {
		_, err := fmt.Fprintf(tab, "%s\t%s\t%s\t%s\t", e[i].Timestamp.Format(time.RFC3339), *e[i].LogicalResourceId, *e[i].ResourceType, *e[i].ResourceStatus)
		if err != nil {
			return err
		}
		fmt.Fprintln(tab)
	}
	tab.Flush()
	return nil
}

func NewCloudformationEventsTailer(stackName string, nbEvents int, enableRefresh bool, frequency time.Duration, watchDeployment bool) *stackEventTailer {
	return &stackEventTailer{stackName: stackName, refresh: enableRefresh, refreshFrequency: frequency, nbEvents: nbEvents, watchDeployment: watchDeployment}
}

func (t *stackEventTailer) Name() string {
	return "stack-events"
}

func (t *stackEventTailer) Tail(w io.Writer) error {
	cfn, ok := awsservices.CloudformationService.(*awsservices.Cloudformation)
	if !ok {
		return fmt.Errorf("invalid cloud service, expected awsservices.Cloudformation, got %T", awsservices.CloudformationService)
	}

	if t.refreshFrequency < 2*time.Second {
		return fmt.Errorf("invalid refresh frequency: %s", t.refreshFrequency)
	}

	if t.watchDeployment {
		deploying, err := t.isStackBeingDeployed(cfn)
		if err != nil {
			return err
		}

		if !deploying {
			return fmt.Errorf("Stack %s not being deployed at the moment", t.stackName)
		}

	} else {
		if err := t.displayLastEvents(cfn, w); err != nil || !t.refresh {
			return err
		}
	}

	ticker := time.NewTicker(t.refreshFrequency)
	defer ticker.Stop()
	for range ticker.C {
		if t.watchDeployment {
			isDeploymentFinished, err := t.displayRelevantEvents(cfn, w)
			if err != nil || isDeploymentFinished {
				return err
			}
		}

		err := t.displayNewEvents(cfn, w)
		if err != nil {
			return err
		}
	}

	return nil
}

// get last N events
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

		stEvents = append(stEvents, resp.StackEvents...)
		if len(stEvents) > t.nbEvents || resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return stEvents, nil
}

// get all events never than last seen event
func (t *stackEventTailer) getNewEvents(cfn *awsservices.Cloudformation) (stackEvents, error) {
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
			if t.lastEventID != nil && *e.EventId == *t.lastEventID {
				break
			}
			stEvents = append(stEvents, e)
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return stEvents, nil
}

func (t *stackEventTailer) displayLastEvents(cfn *awsservices.Cloudformation, w io.Writer) error {
	events, err := t.getLatestEvents(cfn)
	if err != nil {
		return err
	}

	if len(events) > 0 {
		t.lastEventID = events[0].EventId
		return events[:t.nbEvents].printReverse(w)
	}

	return nil
}

func (t *stackEventTailer) displayNewEvents(cfn *awsservices.Cloudformation, w io.Writer) error {
	events, err := t.getNewEvents(cfn)
	if err != nil {
		return err
	}

	if len(events) > 0 {
		t.lastEventID = events[0].EventId
	}

	return events.printReverse(w)
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
			if t.lastEventID != nil && *e.EventId == *t.lastEventID {
				return stEvents, false, nil
			}
			stEvents = append(stEvents, e)
			if *e.ResourceType == "AWS::CloudFormation::Stack" && strings.HasSuffix(*e.ResourceStatus, "_IN_PROGRESS") {
				return stEvents, false, nil
			}

			if *e.ResourceType == "AWS::CloudFormation::Stack" && strings.HasSuffix(*e.ResourceStatus, "_COMPLETE") {
				return stEvents, true, nil
			}

			if *e.ResourceType == "AWS::CloudFormation::Stack" && strings.HasSuffix(*e.ResourceStatus, "_FAILED") {
				return stEvents, true, fmt.Errorf("Deployment failed")
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
