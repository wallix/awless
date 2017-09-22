package awstailers

import (
	"fmt"
	"io"
	"sort"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/wallix/awless/aws/services"
)

type scalingActivitiesTailer struct {
	follow           bool
	pollingFrequency time.Duration
	lastEventTime    time.Time
	nbEvents         int
}

func NewScalingActivitiesTailer(nbEvents int, follow bool, frequency time.Duration) *scalingActivitiesTailer {
	return &scalingActivitiesTailer{nbEvents: nbEvents, follow: follow, pollingFrequency: frequency}
}

func (t *scalingActivitiesTailer) Name() string {
	return "scaling-activities"
}

func (t *scalingActivitiesTailer) Tail(w io.Writer) error {
	infra, ok := awsservices.InfraService.(*awsservices.Infra)
	if !ok {
		return fmt.Errorf("invalid cloud service, expected awsservices.Infra, got %T", awsservices.InfraService)
	}
	if err := t.displayLastEvents(infra, w); err != nil {
		return err
	}

	if t.lastEventTime.IsZero() {
		return nil
	}

	if !t.follow {
		return nil
	}

	if t.pollingFrequency < 5*time.Second {
		return fmt.Errorf("invalid polling frequency: %s", t.pollingFrequency)
	}

	ticker := time.NewTicker(t.pollingFrequency)
	defer ticker.Stop()
	for range ticker.C {
		if err := t.displayNewEvents(infra, w); err != nil {
			return err
		}
	}
	return nil

}

func (t *scalingActivitiesTailer) displayLastEvents(infra *awsservices.Infra, w io.Writer) error {
	out, err := infra.AutoScalingAPI.DescribeScalingActivities(&autoscaling.DescribeScalingActivitiesInput{MaxRecords: awssdk.Int64(int64(t.nbEvents))})
	if err != nil {
		return err
	}
	var events []*event
	for i, activity := range out.Activities {
		evt := newEventFromScalingActivity(activity)
		if i == 0 {
			t.lastEventTime = evt.stamp
		}
		events = append(events, evt)
	}
	sort.Slice(events, func(i int, j int) bool { return events[i].stamp.Before(events[j].stamp) })
	for _, evt := range events {
		if err := evt.print(w); err != nil {
			return err
		}
	}
	return nil
}

func (t *scalingActivitiesTailer) displayNewEvents(infra *awsservices.Infra, w io.Writer) error {
	var eventFound bool
	var newEvents []*event
	lastEventTime := t.lastEventTime
	err := infra.AutoScalingAPI.DescribeScalingActivitiesPages(&autoscaling.DescribeScalingActivitiesInput{}, func(page *autoscaling.DescribeScalingActivitiesOutput, lastPage bool) bool {
		for _, act := range page.Activities {
			evt := newEventFromScalingActivity(act)
			if t.lastEventTime.Before(evt.stamp) {
				t.lastEventTime = evt.stamp
			}
			if evt.stamp == lastEventTime || evt.stamp.Before(lastEventTime) {
				eventFound = true
				break
			}
			newEvents = append(newEvents, evt)
		}
		return !eventFound
	})
	if err != nil {
		return err
	}
	sort.Slice(newEvents, func(i int, j int) bool { return newEvents[i].stamp.Before(newEvents[j].stamp) })
	for _, e := range newEvents {
		if err := e.print(w); err != nil {
			return err
		}
	}
	return nil
}

type event struct {
	id      string
	element string
	stamp   time.Time
	message string
}

func newEventFromScalingActivity(s *autoscaling.Activity) *event {
	return &event{
		id:      awssdk.StringValue(s.ActivityId),
		stamp:   awssdk.TimeValue(s.StartTime),
		message: fmt.Sprintf("%s: %s", awssdk.StringValue(s.StatusCode), awssdk.StringValue(s.Description)),
		element: awssdk.StringValue(s.AutoScalingGroupName),
	}
}

func (e *event) print(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s: %s\n\t%s\n", e.stamp, e.element, e.message)
	return err
}
