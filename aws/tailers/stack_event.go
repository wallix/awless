package awstailers

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
	"github.com/wallix/awless/aws/services"
)

const (
	FilterStackEventLogicalID    = "id"
	FilterStackEventTimestamp    = "ts"
	FilterStackEventStatus       = "status"
	FilterStackEventStatusReason = "reason"
	FilterStackEventType         = "type"

	// valid stack status codes
	// http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-describing-stacks.html#w2ab2c15c15c17c11
	cfStackEventUpdateInProgress = "UPDATE_IN_PROGRESS"
	cfStackEventCreateInProgress = "CREATE_IN_PROGRESS"
	cfStackEventDeleteInProgress = "DELETE_IN_PROGRESS"
	cfStackEventCompleteSuffix   = "_COMPLETE"
	cfStackEventFailedSuffix     = "_FAILED"
	cfStackEventInProgressSuffix = "_IN_PROGRESS"
	cfStackType                  = "AWS::CloudFormation::Stack"
)

type filters []string

type stackEventTailer struct {
	stackName        string
	follow           bool
	pollingFrequency time.Duration
	lastEventID      *string
	nbEvents         int
	filters          filters
	deploymentStatus deploymentStatus
}

type stackEvent struct {
	*cloudformation.StackEvent
}

type stackEvents []stackEvent

func NewCloudformationEventsTailer(stackName string, nbEvents int, enableFollow bool, frequency time.Duration, f filters) *stackEventTailer {
	return &stackEventTailer{
		stackName:        stackName,
		follow:           enableFollow,
		pollingFrequency: frequency,
		nbEvents:         nbEvents,
		filters:          f,
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

	// create new tabwriter and
	// add header based on filters
	tab := tabwriter.NewWriter(w, 25, 8, 0, '\t', 0)
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
	defer ticker.Stop()
	for range ticker.C {
		err := t.displayRelevantEvents(cfn, tab)
		if err != nil {
			return err
		}

		tab.Flush()

		if t.deploymentStatus.isFinished {
			if len(t.deploymentStatus.failedEvents) > 0 {
				var errBuf bytes.Buffer
				var f filters = []string{FilterStackEventLogicalID, FilterStackEventType, FilterStackEventStatus, FilterStackEventStatusReason}

				errBuf.WriteString("Deployment failed.\nFailed events summary:\n")
				errBuf.Write(f.header())

				t.deploymentStatus.failedEvents.printReverse(&errBuf, f)

				return fmt.Errorf(errBuf.String())
			}

			return nil
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
	case strings.HasSuffix(str, cfStackEventInProgressSuffix):
		return color.New(color.FgYellow).SprintFunc()(str)
	case strings.HasSuffix(str, cfStackEventCompleteSuffix):
		return color.New(color.FgGreen).SprintFunc()(str)
	case strings.HasSuffix(str, cfStackEventFailedSuffix):
		return color.New(color.FgRed).SprintFunc()(str)
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
	// var bold = color.New(color.Bo)
	for i, filter := range f {
		switch filter {
		case FilterStackEventLogicalID:
			buf.WriteString("Logical ID")
		case FilterStackEventTimestamp:
			buf.WriteString("Timestamp")
		case FilterStackEventStatus:
			buf.WriteString("Status")
		case FilterStackEventStatusReason:
			buf.WriteString("Status Reason")
		case FilterStackEventType:
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
		case f == FilterStackEventLogicalID && e.LogicalResourceId != nil:
			buf.WriteString(*e.LogicalResourceId)
		case f == FilterStackEventTimestamp && e.Timestamp != nil:
			buf.WriteString(e.Timestamp.Format(time.RFC3339))
		case f == FilterStackEventStatus && e.ResourceStatus != nil:
			buf.WriteString(coloredResourceStatus(*e.ResourceStatus))
		case f == FilterStackEventStatusReason && e.ResourceStatusReason != nil:
			buf.WriteString(*e.ResourceStatusReason)
		case f == FilterStackEventType && e.ResourceType != nil:
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
	return (s.ResourceType != nil && *s.ResourceType == cfStackType) &&
		(s.ResourceStatus != nil &&
			*s.ResourceStatus == cfStackEventCreateInProgress ||
			*s.ResourceStatus == cfStackEventDeleteInProgress ||
			*s.ResourceStatus == cfStackEventUpdateInProgress)
}

func (s *stackEvent) isDeploymentFinished() bool {
	return (s.ResourceType != nil && *s.ResourceType == cfStackType) &&
		(s.ResourceStatus != nil &&
			strings.HasSuffix(*s.ResourceStatus, cfStackEventCompleteSuffix) ||
			strings.HasSuffix(*s.ResourceStatus, cfStackEventFailedSuffix))
}

func (s *stackEvent) isFailed() bool {
	return (s.ResourceStatus != nil && strings.HasSuffix(*s.ResourceStatus, cfStackEventFailedSuffix))
}
