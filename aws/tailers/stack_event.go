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

func (e *stackEvents) print(w io.Writer) error {
	tab := tabwriter.NewWriter(w, 35, 8, 0, '\t', 0)
	for _, event := range *e {
		_, err := fmt.Fprintf(tab, "%s\t%s\t%s\t%s\t", event.Timestamp.Format(time.RFC3339), *event.LogicalResourceId, *event.ResourceType, *event.ResourceStatus)
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

	if t.refreshFrequency < 5*time.Second {
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

func (t *stackEventTailer) getAllEvents(cfn *awsservices.Cloudformation) ([]*cloudformation.StackEvent, error) {
	params := &cloudformation.DescribeStackEventsInput{
		StackName: &t.stackName,
	}

	var events []*cloudformation.StackEvent

	for {
		resp, err := cfn.DescribeStackEvents(params)
		if err != nil {
			return nil, err
		}

		events = append(events, resp.StackEvents...)
		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return events, nil
}

func (t *stackEventTailer) displayLastEvents(cfn *awsservices.Cloudformation, w io.Writer) error {
	events, err := t.getAllEvents(cfn)
	if err != nil {
		return err
	}
	t.lastEventID = events[len(events)-1].EventId
	lastEvents := stackEvents(events[:t.nbEvents])

	return lastEvents.print(w)
}

func (t *stackEventTailer) displayNewEvents(cfn *awsservices.Cloudformation, w io.Writer) error {
	events, err := t.getAllEvents(cfn)
	if err != nil {
		return err
	}

	var newEvents stackEvents

	for i, e := range events {
		if e.EventId == t.lastEventID {
			newEvents = stackEvents(events[i:])
			break
		}
	}

	if len(newEvents) > 0 {
		t.lastEventID = events[len(newEvents)-1].EventId
		return newEvents.print(w)
	}

	return nil
}
