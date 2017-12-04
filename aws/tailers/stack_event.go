package awstailers

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/fatih/color"
	"github.com/wallix/awless/aws/services"
)

const (
	StackEventLogicalID    = "id"
	StackEventTimestamp    = "ts"
	StackEventStatus       = "status"
	StackEventStatusReason = "reason"
	StackEventType         = "type"

	// valid stack status codes
	// http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-describing-stacks.html#w2ab2c15c15c17c11
	StackEventComplete   = "COMPLETE"
	StackEventFailed     = "FAILED"
	StackEventInProgress = "IN_PROGRESS"
)

type filters []string

type stackEventTailer struct {
	stackName          string
	follow             bool
	pollingFrequency   time.Duration
	lastEventID        *string
	nbEvents           int
	filters            filters
	deploymentStatus   deploymentStatus
	timeout            time.Duration
	cancelAfterTimeout bool
}

type stackEvent struct {
	*cloudformation.StackEvent
}

type stackEvents []stackEvent

func NewCloudformationEventsTailer(stackName string, nbEvents int, enableFollow bool, frequency time.Duration, f filters, timeout time.Duration, cancelAfterTimeout bool) *stackEventTailer {
	return &stackEventTailer{
		stackName:          stackName,
		follow:             enableFollow,
		pollingFrequency:   frequency,
		nbEvents:           nbEvents,
		filters:            f,
		timeout:            timeout,
		cancelAfterTimeout: cancelAfterTimeout,
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
		return fmt.Errorf("invalid polling frequency: %s, must be greater than 5s", t.pollingFrequency)
	}

	tab := tabwriter.NewWriter(w, 8, 8, 8, '\t', 0)
	tab.Write(t.filters.header())

	if !t.follow {
		if err := t.displayLastEvents(cfn, tab); err != nil {
			return err
		}

		tab.Flush()

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
	timer := time.NewTimer(t.timeout)

	defer ticker.Stop()
	defer timer.Stop()

	isTimeoutReached := false
	for {
		select {
		case <-timer.C:
			isTimeoutReached = true
			if t.cancelAfterTimeout {
				color.Red("Timeout (%s) reached.", t.timeout.String())
				color.Red("Canceling update of stack %q", t.stackName)
				err := t.cancelStackUpdate(cfn)
				if err != nil {
					return fmt.Errorf("Couldn't cancel stack update.\nError: %s\nStack update could be running, please check manually", err)
				}
			} else {
				return fmt.Errorf("Timeout (%s) reached. Exiting...", t.timeout.String())
			}
		case <-ticker.C:
			if err := t.displayRelevantEvents(cfn, tab); err != nil {
				return err
			}

			tab.Flush()

			if t.deploymentStatus.isFinished {
				if len(t.deploymentStatus.failedEvents) > 0 {
					var errBuf bytes.Buffer
					var f filters = []string{StackEventLogicalID, StackEventType, StackEventStatus, StackEventStatusReason}

					if isTimeoutReached {
						errBuf.WriteString("Update was cancelled because timeout has been reached and option 'Cancel On Timeout' enabled\n")
					} else {
						errBuf.WriteString("Update failed\n")
					}

					errBuf.WriteString("Failed events summary:\n")

					// printing error events as a nice table
					errTab := tabwriter.NewWriter(&errBuf, 25, 8, 0, '\t', 0)
					errTab.Write(f.header())
					t.deploymentStatus.failedEvents.printReverse(errTab, f)
					errTab.Flush()

					return fmt.Errorf(errBuf.String())
				}
				return nil
			}
		}
	}
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
			stEvents = append(stEvents, stackEvent{e})
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
		return events.printReverse(w, t.filters)
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

	return strings.HasSuffix(*stacks.Stacks[0].StackStatus, StackEventInProgress), nil
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
			event := stackEvent{e}
			// if lastEventID == nil then it's first run of this method
			// if lastEventID == nil then it's not first run and print only new messages
			if t.lastEventID != nil && *e.EventId == *t.lastEventID {
				return stEvents, nil
			}
			stEvents = append(stEvents, event)

			// looking for the message which says that stack update or create started
			// making it as a first messages in the deployment events
			if event.isDeploymentStart() {
				return stEvents, nil
			}

			// if we found message, that stack create/update/delete completed or failed
			// then marking build as complete, but keep tailing
			if event.isDeploymentFinished() {
				t.deploymentStatus.isFinished = true
			}

			// if we found fail message then append error to the slice
			// but keep tailing
			if event.isFailed() {
				t.deploymentStatus.failedEvents = append(t.deploymentStatus.failedEvents, event)
			}
		}

		if resp.NextToken == nil {
			return stEvents, nil
		}

		params.NextToken = resp.NextToken
	}

}

func (t *stackEventTailer) displayRelevantEvents(cfn *awsservices.Cloudformation, w io.Writer) error {
	events, err := t.getRelevantEvents(cfn)
	if err != nil {
		return err
	}

	if len(events) > 0 {
		t.lastEventID = events[0].EventId
	}

	return events.printReverse(w, t.filters)
}

func coloredResourceStatus(str string) string {
	switch {
	case strings.HasSuffix(str, StackEventFailed),
		str == cloudformation.StackStatusUpdateRollbackInProgress,
		str == cloudformation.StackStatusRollbackInProgress:
		return color.New(color.FgRed).SprintFunc()(str)
	case strings.HasSuffix(str, StackEventInProgress):
		return color.New(color.FgYellow).SprintFunc()(str)
	case strings.HasSuffix(str, StackEventComplete):
		return color.New(color.FgGreen).SprintFunc()(str)
	default:
		return str
	}

}

func (e stackEvents) printReverse(w io.Writer, f filters) error {
	for i := len(e) - 1; i >= 0; i-- {
		w.Write(e[i].filter(f))
	}

	return nil
}

func (f filters) header() []byte {
	var buf bytes.Buffer
	for i, filter := range f {
		switch filter {
		case StackEventLogicalID:
			buf.WriteString("Logical ID")
		case StackEventTimestamp:
			buf.WriteString("Timestamp")
		case StackEventStatus:
			buf.WriteString("Status")
		case StackEventStatusReason:
			buf.WriteString("Status Reason")
		case StackEventType:
			buf.WriteString("Type")
		}

		if i != len(f)-1 {
			buf.WriteRune('\t')
		}

	}

	// with "\n" formatted with bold, tabwriter somehow shift lines
	// so we need to add "\n" after string being bolded
	return []byte(color.New(color.Bold).Sprintf(buf.String()) + "\n")
}

func (e *stackEvent) filter(filters []string) (out []byte) {
	var buf bytes.Buffer

	for i, f := range filters {
		switch {
		case f == StackEventLogicalID && e.LogicalResourceId != nil:
			buf.WriteString(*e.LogicalResourceId)
		case f == StackEventTimestamp && e.Timestamp != nil:
			buf.WriteString(e.Timestamp.Format(time.RFC3339))
		case f == StackEventStatus && e.ResourceStatus != nil:
			buf.WriteString(coloredResourceStatus(*e.ResourceStatus))
		case f == StackEventStatusReason && e.ResourceStatusReason != nil:
			buf.WriteString(*e.ResourceStatusReason)
		case f == StackEventType && e.ResourceType != nil:
			buf.WriteString(*e.ResourceType)
		}

		if i != len(filters)-1 {
			buf.WriteRune('\t')
		}

	}

	buf.WriteRune('\n')

	return buf.Bytes()
}

func (s *stackEvent) isDeploymentStart() bool {
	return (s.ResourceType != nil && *s.ResourceType == configservice.ResourceTypeAwsCloudFormationStack) &&
		(s.ResourceStatus != nil &&
			*s.ResourceStatus == cloudformation.ResourceStatusCreateInProgress ||
			*s.ResourceStatus == cloudformation.ResourceStatusDeleteInProgress ||
			*s.ResourceStatus == cloudformation.ResourceStatusUpdateInProgress)
}

func (s *stackEvent) isDeploymentFinished() bool {
	return (s.ResourceType != nil && *s.ResourceType == configservice.ResourceTypeAwsCloudFormationStack) &&
		(s.ResourceStatus != nil &&
			strings.HasSuffix(*s.ResourceStatus, StackEventComplete) ||
			strings.HasSuffix(*s.ResourceStatus, StackEventFailed))
}

func (s *stackEvent) isFailed() bool {
	return (s.ResourceStatus != nil && (strings.HasSuffix(*s.ResourceStatus, StackEventFailed) || *s.ResourceStatus == cloudformation.StackStatusUpdateRollbackInProgress))
}

func (s *stackEventTailer) cancelStackUpdate(cfn *awsservices.Cloudformation) error {
	inp := &cloudformation.CancelUpdateStackInput{StackName: &s.stackName}
	_, err := cfn.CancelUpdateStack(inp)
	return err
}
