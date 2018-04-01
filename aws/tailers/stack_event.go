package awstailers

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/fatih/color"
	"github.com/wallix/awless/aws/services"
)

const (
	StackEventFilterLogicalID    = "id"
	StackEventFilterTimestamp    = "ts"
	StackEventFilterStatus       = "status"
	StackEventFilterStatusReason = "reason"
	StackEventFilterType         = "type"
	StackEventFilterPhysicalId   = "physical-id"

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

// Copy of cloudformation.StackEvent for futher string formating
type stackEvent struct {
	Timestamp         *string `width:"20,5"`
	ResourceStatus    *string `width:"50,5"`
	ResourceType      *string `width:"45,5"`
	LogicalResourceId *string `width:"20,5"`

	PhysicalResourceId   *string `width:"50,5"`
	ResourceStatusReason *string `width:"50,5"`
	EventId              *string
}

var filtersMapping = map[string]string{
	StackEventFilterLogicalID:    "LogicalResourceId",
	StackEventFilterTimestamp:    "Timestamp",
	StackEventFilterStatus:       "ResourceStatus",
	StackEventFilterStatusReason: "ResourceStatusReason",
	StackEventFilterType:         "ResourceType",
	StackEventFilterPhysicalId:   "PhysicalResourceId",
}

var DefaultStackEventFilters = []string{StackEventFilterTimestamp, StackEventFilterLogicalID, StackEventFilterType, StackEventFilterStatus}
var AllStackEventFilters = append(DefaultStackEventFilters, StackEventFilterStatusReason, StackEventFilterPhysicalId)

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

	w.Write(t.filters.header())
	if !t.follow {
		return t.displayLastEvents(cfn, w)
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
			if err := t.displayRelevantEvents(cfn, w); err != nil {
				return err
			}

			if t.deploymentStatus.isFinished {
				if len(t.deploymentStatus.failedEvents) > 0 {
					var errBuf bytes.Buffer
					var f filters = []string{StackEventFilterLogicalID, StackEventFilterType, StackEventFilterStatus, StackEventFilterStatusReason}

					if isTimeoutReached {
						errBuf.WriteString("Update was cancelled because timeout has been reached and option 'Cancel On Timeout' enabled\n")
					} else {
						errBuf.WriteString("Update failed\n")
					}

					errBuf.WriteString("Failed events summary:\n")

					// using tabwriter here, because we have all data
					// and no need to stream it
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
			stEvents = append(stEvents, NewStackEvent(e))
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
			event := NewStackEvent(e)
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

func colorizeResourceStatus(str string) *string {
	var c color.Attribute
	switch {
	case strings.HasSuffix(str, StackEventFailed),
		str == cloudformation.StackStatusUpdateRollbackInProgress,
		str == cloudformation.StackStatusRollbackInProgress:
		c = color.FgRed
	case strings.HasSuffix(str, StackEventInProgress):
		c = color.FgYellow
	case strings.HasSuffix(str, StackEventComplete):
		c = color.FgGreen
	}

	s := color.New(c).SprintFunc()(str)

	return &s
}

func (e stackEvents) printReverse(w io.Writer, f filters) error {
	for i := len(e) - 1; i >= 0; i-- {
		w.Write(e[i].format(f))
	}

	return nil
}

func (f filters) header() []byte {
	s := &stackEvent{
		Timestamp:            func() *string { t := color.New(color.Bold).Sprintf("Timestamp"); return &t }(),
		ResourceStatus:       func() *string { t := color.New(color.Bold).Sprintf("Status"); return &t }(),
		LogicalResourceId:    func() *string { t := color.New(color.Bold).Sprintf("Logical ID"); return &t }(),
		PhysicalResourceId:   func() *string { t := color.New(color.Bold).Sprintf("Physical ID"); return &t }(),
		ResourceStatusReason: func() *string { t := color.New(color.Bold).Sprintf("Status Reason"); return &t }(),
		ResourceType:         func() *string { t := color.New(color.Bold).Sprintf("Type"); return &t }(),
	}

	return s.format(f)
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

func (s *stackEvent) fromCFEvent() bool {
	return (s.ResourceStatus != nil && (strings.HasSuffix(*s.ResourceStatus, StackEventFailed) || *s.ResourceStatus == cloudformation.StackStatusUpdateRollbackInProgress))
}

func (s *stackEventTailer) cancelStackUpdate(cfn *awsservices.Cloudformation) error {
	inp := &cloudformation.CancelUpdateStackInput{StackName: &s.stackName}
	_, err := cfn.CancelUpdateStack(inp)
	return err
}

func NewStackEvent(e *cloudformation.StackEvent) stackEvent {
	return stackEvent{
		Timestamp:            func() *string { t := e.Timestamp.Format(time.RFC3339); return &t }(),
		ResourceStatus:       colorizeResourceStatus(*e.ResourceStatus),
		ResourceType:         e.ResourceType,
		LogicalResourceId:    e.LogicalResourceId,
		PhysicalResourceId:   e.PhysicalResourceId,
		ResourceStatusReason: e.ResourceStatusReason,
		EventId:              e.EventId,
	}
}

// Format reads the struct tag `width:"<width>,<space>"`
// further marshaling into structured field
func (s *stackEvent) format(fil filters) []byte {
	tp := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()

	fmt.Println(" Call Format")

	buf := bytes.Buffer{}
	var nextLine *stackEvent
	for _, f := range fil {
		field, ok := tp.FieldByName(filtersMapping[f])
		if !ok {
			continue
		}
		value := v.FieldByName(filtersMapping[f])

		splt := strings.Split(field.Tag.Get("width"), ",")
		if len(splt) != 2 {
			continue
		}

		width, err := strconv.Atoi(splt[0])
		if err != nil {
			continue
		}

		space, err := strconv.Atoi(splt[1])
		if err != nil {
			continue
		}

		var v string
		if !value.IsNil() {
			v = value.Elem().String()
		}

		// handle coloring
		// if string starts with "\x1b" then it is colored
		if strings.HasPrefix(v, "\x1b") {
			// color adds additional length to the string
			// which is not displayed in the console
			// and results in text shift
			// so we need to increase column width a bit
			// colored string looks like: "\x1b[31mText\x1b[0m"
			width += strings.Index(v, "m") + 1 + len("\x1b[0m")
		}

		if len(v) > width {
			if nextLine == nil {
				nextLine = &stackEvent{}
			}
			if field.Name == "ResourceStatusReason" {
				nextLine.ResourceStatusReason = aws.String(v[width:])
			}
			if field.Name == "ResourceStatus" {
				nextLine.ResourceStatus = aws.String(v[width:])
			}

			v = v[:width]
		}

		buf.WriteString(v)
		// fil the rest of the line space with " "
		buf.WriteString(createSpaces(width + space - len(v)))
		// totalWidth += width
	}

	buf.WriteRune('\n')
	if nextLine != nil {
		buf.Write(nextLine.format(fil))
	}
	return buf.Bytes()
}

func createSpaces(n int) string {
	var buf = bytes.Buffer{}
	for i := 0; i < n; i++ {
		buf.WriteString(".")
	}

	return buf.String()
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
