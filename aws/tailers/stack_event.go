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
	filters          []string
	deploymentStatus deploymentStatus
}

type stackEvents []*cloudformation.StackEvent

func NewCloudformationEventsTailer(stackName string, nbEvents int, enableFollow bool, frequency time.Duration) *stackEventTailer {
	return &stackEventTailer{
		stackName:        stackName,
		follow:           enableFollow,
		pollingFrequency: frequency,
		nbEvents:         nbEvents,
		filters:          []string{"Timestamp", "Logical ID", "Type", "Status"}, // TODO: allow to set filters from command line
	}
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

	// if this is first run, print header to the stack events output
	var isFirstRun = true

	ticker := time.NewTicker(t.pollingFrequency)
	defer ticker.Stop()
	for range ticker.C {
		err := t.displayRelevantEvents(cfn, w, isFirstRun)
		if err != nil {
			return err
		}

		isFirstRun = false

		if len(t.deploymentStatus.failedEvents) > 0 {
			t.deploymentStatus.failedEvents.printErrorsReverse(w)
			return fmt.Errorf("Deployment failed")
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
		return events.printReverse(w, t.filters, true)
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

type deploymentStatus struct {
	isFinished   bool
	failedEvents stackEvents
}

// get last N events relevant for current deployment in progress
func (t *stackEventTailer) getRelevantEvents(cfn *awsservices.Cloudformation) (stEvents stackEvents, err error) {
	params := &cloudformation.DescribeStackEventsInput{
		StackName: &t.stackName,
	}

	var resp *cloudformation.DescribeStackEventsOutput

	for {
		resp, err = cfn.DescribeStackEvents(params)
		if err != nil {
			return nil, err
		}

		for _, e := range resp.StackEvents {
			// if lastEventID == nil then it's first run of this method
			// if lastEventID == nil then it's not first run and print only new messages
			if t.lastEventID != nil && *e.EventId == *t.lastEventID {
				return stEvents, nil
			}
			stEvents = append(stEvents, e)

			// looking for the message which says that stack update or create started
			// making it as a first messages in the deployment events
			if *e.ResourceType == "AWS::CloudFormation::Stack" &&
				(strings.HasSuffix(*e.ResourceStatus, "UPDATE_IN_PROGRESS") || strings.HasSuffix(*e.ResourceStatus, "CREATE_IN_PROGRESS")) {
				return stEvents, nil
			}

			// if we found message, that stack create/update completed
			// then marking build as complete, but keep tailing
			if *e.ResourceType == "AWS::CloudFormation::Stack" && strings.HasSuffix(*e.ResourceStatus, "_COMPLETE") {
				t.deploymentStatus.isFinished = true
			}

			// if we found message, that stack create/update failed
			// then marking build as complete, but keep tailing
			if *e.ResourceType == "AWS::CloudFormation::Stack" && strings.HasSuffix(*e.ResourceStatus, "_FAILED") {
				t.deploymentStatus.isFinished = true
			}

			// if we found next message, then any resource create/update failed,
			// then marking build as failed, but keep tailing
			if strings.HasSuffix(*e.ResourceStatus, "_FAILED") {
				t.deploymentStatus.failedEvents = append(t.deploymentStatus.failedEvents, e)
			}
		}

		if resp.NextToken == nil {
			return stEvents, nil
		}

		params.NextToken = resp.NextToken
	}

}

func (t *stackEventTailer) displayRelevantEvents(cfn *awsservices.Cloudformation, w io.Writer, isFirstRun bool) error {
	events, err := t.getRelevantEvents(cfn)
	if err != nil {
		return err
	}

	if len(events) > 0 {
		t.lastEventID = events[0].EventId
	}

	return events.printReverse(w, t.filters, isFirstRun)
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
func (e stackEvents) printReverse(w io.Writer, filters []string, withHeader bool) error {
	tab := tabwriter.NewWriter(w, 25, 8, 0, '\t', 0)
	if withHeader {
		tab.Write([]byte(color.New(color.Bold).Sprint(strings.Join(filters, "\t"))))
		tab.Write([]byte("\n"))
	}
	for i := len(e) - 1; i >= 0; i-- {
		// TODO: Create output based on filters
		_, err := fmt.Fprintf(tab, "%s\t%s\t%s\t%s\n", e[i].Timestamp.Format(time.RFC3339), *e[i].LogicalResourceId, *e[i].ResourceType, coloredResourceStatus(*e[i].ResourceStatus))
		if err != nil {
			return err
		}
	}
	tab.Flush()
	return nil
}

func (e stackEvents) printErrorsReverse(w io.Writer) error {
	tab := tabwriter.NewWriter(w, 25, 8, 0, '\t', 0)

	tab.Write([]byte("\nFailed events summary:\n"))
	tab.Write([]byte(color.New(color.Bold).Sprint("Logical ID\tType\tStatus\tStatus Reason\n")))

	for i := len(e) - 1; i >= 0; i-- {
		_, err := fmt.Fprintf(tab, "%s\t%s\t%s\t%s\n", *e[i].LogicalResourceId, *e[i].ResourceType, coloredResourceStatus(*e[i].ResourceStatus), *e[i].ResourceStatusReason)
		if err != nil {
			return err
		}
	}
	tab.Flush()
	return nil
}
