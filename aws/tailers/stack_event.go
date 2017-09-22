package awstailers

import (
	"fmt"
	"io"
	"time"

	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/wallix/awless/aws/services"
)

type stackEventTailer struct {
	stackName        string
	refresh          bool
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

func NewCloudformationEventsTailer(stackName string, nbEvents int, enableRefresh bool, frequency time.Duration) *stackEventTailer {
	return &stackEventTailer{stackName: stackName, refresh: enableRefresh, refreshFrequency: frequency, nbEvents: nbEvents}
}

func (t *stackEventTailer) Name() string {
	return "stack-events"
}

func (t *stackEventTailer) Tail(w io.Writer) error {
	cfn, ok := awsservices.CloudformationService.(*awsservices.Cloudformation)
	if !ok {
		return fmt.Errorf("invalid cloud service, expected awsservices.Cloudformation, got %T", awsservices.CloudformationService)
	}

	if err := t.displayLastEvents(cfn, w); err != nil || !t.refresh {
		return err
	}

	if t.refreshFrequency < 2*time.Second {
		return fmt.Errorf("invalid refresh frequency: %s", t.refreshFrequency)
	}

	ticker := time.NewTicker(t.refreshFrequency)
	defer ticker.Stop()
	for range ticker.C {
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
			if *e.EventId == *t.lastEventID {
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
