package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func TestAlarm(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create alarm name=my-alarm-name evaluation-periods=3 metric=my-metric-name namespace=my-namespace-name "+
			"operator=GreaterThanThreshold period=4 statistic-function=Average threshold=12.5 "+
			"alarm-actions=arn:my:alarm:action1,arn:my:alarm:action2 description='This is a test alarm' "+
			"dimensions=DimKey1:val1,DimKey2:val2 enabled=true insufficientdata-actions=arn:my:action:insufficientdata "+
			"ok-actions=arn:my:ok-action unit=Bytes").
			Mock(&cloudwatchMock{
				PutMetricAlarmFunc: func(param0 *cloudwatch.PutMetricAlarmInput) (*cloudwatch.PutMetricAlarmOutput, error) {
					return &cloudwatch.PutMetricAlarmOutput{}, nil
				},
			}).ExpectInput("PutMetricAlarm", &cloudwatch.PutMetricAlarmInput{
			AlarmName:               String("my-alarm-name"),
			EvaluationPeriods:       Int64(3),
			MetricName:              String("my-metric-name"),
			Namespace:               String("my-namespace-name"),
			ComparisonOperator:      String("GreaterThanThreshold"),
			Period:                  Int64(4),
			Statistic:               String("Average"),
			Threshold:               Float64(12.5),
			AlarmActions:            []*string{String("arn:my:alarm:action1"), String("arn:my:alarm:action2")},
			AlarmDescription:        String("This is a test alarm"),
			Dimensions:              []*cloudwatch.Dimension{{Name: String("DimKey1"), Value: String("val1")}, {Name: String("DimKey2"), Value: String("val2")}},
			ActionsEnabled:          Bool(true),
			InsufficientDataActions: []*string{String("arn:my:action:insufficientdata")},
			OKActions:               []*string{String("arn:my:ok-action")},
			Unit:                    String("Bytes"),
		}).ExpectCommandResult("my-alarm-name").ExpectCalls("PutMetricAlarm").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete alarm name=my-alarm-to-delete").Mock(&cloudwatchMock{
			DeleteAlarmsFunc: func(param0 *cloudwatch.DeleteAlarmsInput) (*cloudwatch.DeleteAlarmsOutput, error) { return nil, nil },
		}).ExpectInput("DeleteAlarms", &cloudwatch.DeleteAlarmsInput{AlarmNames: []*string{String("my-alarm-to-delete")}}).
			ExpectCalls("DeleteAlarms").Run(t)

		Template("delete alarm name=alarm1,alarm2").Mock(&cloudwatchMock{
			DeleteAlarmsFunc: func(param0 *cloudwatch.DeleteAlarmsInput) (*cloudwatch.DeleteAlarmsOutput, error) { return nil, nil },
		}).ExpectInput("DeleteAlarms", &cloudwatch.DeleteAlarmsInput{AlarmNames: []*string{String("alarm1"), String("alarm2")}}).
			ExpectCalls("DeleteAlarms").Run(t)
	})

	t.Run("start", func(t *testing.T) {
		Template("start alarm names=my-alarm-to-start").Mock(&cloudwatchMock{
			EnableAlarmActionsFunc: func(param0 *cloudwatch.EnableAlarmActionsInput) (*cloudwatch.EnableAlarmActionsOutput, error) {
				return nil, nil
			},
		}).ExpectInput("EnableAlarmActions", &cloudwatch.EnableAlarmActionsInput{AlarmNames: []*string{String("my-alarm-to-start")}}).
			ExpectCalls("EnableAlarmActions").Run(t)

		Template("start alarm names=alarm1,alarm2").Mock(&cloudwatchMock{
			EnableAlarmActionsFunc: func(param0 *cloudwatch.EnableAlarmActionsInput) (*cloudwatch.EnableAlarmActionsOutput, error) {
				return nil, nil
			},
		}).ExpectInput("EnableAlarmActions", &cloudwatch.EnableAlarmActionsInput{AlarmNames: []*string{String("alarm1"), String("alarm2")}}).
			ExpectCalls("EnableAlarmActions").Run(t)
	})

	t.Run("stop", func(t *testing.T) {
		Template("stop alarm names=my-alarm-to-stop").Mock(&cloudwatchMock{
			DisableAlarmActionsFunc: func(param0 *cloudwatch.DisableAlarmActionsInput) (*cloudwatch.DisableAlarmActionsOutput, error) {
				return nil, nil
			},
		}).ExpectInput("DisableAlarmActions", &cloudwatch.DisableAlarmActionsInput{AlarmNames: []*string{String("my-alarm-to-stop")}}).
			ExpectCalls("DisableAlarmActions").Run(t)

		Template("stop alarm names=alarm1,alarm2").Mock(&cloudwatchMock{
			DisableAlarmActionsFunc: func(param0 *cloudwatch.DisableAlarmActionsInput) (*cloudwatch.DisableAlarmActionsOutput, error) {
				return nil, nil
			},
		}).ExpectInput("DisableAlarmActions", &cloudwatch.DisableAlarmActionsInput{AlarmNames: []*string{String("alarm1"), String("alarm2")}}).
			ExpectCalls("DisableAlarmActions").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach alarm name=my-alarm-to-attach action-arn=arn:of:new_action").Mock(&cloudwatchMock{
			DescribeAlarmsFunc: func(param0 *cloudwatch.DescribeAlarmsInput) (*cloudwatch.DescribeAlarmsOutput, error) {
				return &cloudwatch.DescribeAlarmsOutput{
					MetricAlarms: []*cloudwatch.MetricAlarm{{
						ActionsEnabled:                   Bool(true),
						AlarmActions:                     []*string{String("old_action_1"), String("old_action_2")},
						AlarmDescription:                 String("my description"),
						AlarmName:                        String("my-alarm-name"),
						ComparisonOperator:               String("GreaterThanOrEqualToThreshold"),
						Dimensions:                       []*cloudwatch.Dimension{{Name: String("Dim1"), Value: String("Val1")}, {Name: String("Dim2"), Value: String("Val2")}},
						EvaluateLowSampleCountPercentile: String("evaluate"),
						EvaluationPeriods:                Int64(42),
						ExtendedStatistic:                String("p0.5"),
						InsufficientDataActions:          []*string{String("insufic_action_1"), String("insufic_action_2")},
						MetricName:                       String("my-metric-name"),
						Namespace:                        String("my-namespace"),
						OKActions:                        []*string{String("ok_action_1"), String("ok_action_2")},
						Period:                           Int64(17),
						Statistic:                        String("Average"),
						Threshold:                        Float64(42.17),
						TreatMissingData:                 String("ignore"),
						Unit:                             String("Kilobits"),
					}},
				}, nil
			},
			PutMetricAlarmFunc: func(param0 *cloudwatch.PutMetricAlarmInput) (*cloudwatch.PutMetricAlarmOutput, error) {
				return nil, nil
			},
		}).ExpectInput("PutMetricAlarm", &cloudwatch.PutMetricAlarmInput{
			ActionsEnabled:                   Bool(true),
			AlarmActions:                     []*string{String("old_action_1"), String("old_action_2"), String("arn:of:new_action")},
			AlarmDescription:                 String("my description"),
			AlarmName:                        String("my-alarm-name"),
			ComparisonOperator:               String("GreaterThanOrEqualToThreshold"),
			Dimensions:                       []*cloudwatch.Dimension{{Name: String("Dim1"), Value: String("Val1")}, {Name: String("Dim2"), Value: String("Val2")}},
			EvaluateLowSampleCountPercentile: String("evaluate"),
			EvaluationPeriods:                Int64(42),
			ExtendedStatistic:                String("p0.5"),
			InsufficientDataActions:          []*string{String("insufic_action_1"), String("insufic_action_2")},
			MetricName:                       String("my-metric-name"),
			Namespace:                        String("my-namespace"),
			OKActions:                        []*string{String("ok_action_1"), String("ok_action_2")},
			Period:                           Int64(17),
			Statistic:                        String("Average"),
			Threshold:                        Float64(42.17),
			TreatMissingData:                 String("ignore"),
			Unit:                             String("Kilobits"),
		}).ExpectInput("DescribeAlarms", &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []*string{String("my-alarm-to-attach")},
		}).
			ExpectCalls("PutMetricAlarm", "DescribeAlarms").Run(t)
	})
	t.Run("detach", func(t *testing.T) {
		Template("detach alarm name=my-alarm-to-detach action-arn=old_action_2").Mock(&cloudwatchMock{
			DescribeAlarmsFunc: func(param0 *cloudwatch.DescribeAlarmsInput) (*cloudwatch.DescribeAlarmsOutput, error) {
				return &cloudwatch.DescribeAlarmsOutput{
					MetricAlarms: []*cloudwatch.MetricAlarm{{
						ActionsEnabled:                   Bool(true),
						AlarmActions:                     []*string{String("old_action_1"), String("old_action_2")},
						AlarmDescription:                 String("my description"),
						AlarmName:                        String("my-alarm-name"),
						ComparisonOperator:               String("GreaterThanOrEqualToThreshold"),
						Dimensions:                       []*cloudwatch.Dimension{{Name: String("Dim1"), Value: String("Val1")}, {Name: String("Dim2"), Value: String("Val2")}},
						EvaluateLowSampleCountPercentile: String("evaluate"),
						EvaluationPeriods:                Int64(42),
						ExtendedStatistic:                String("p0.5"),
						InsufficientDataActions:          []*string{String("insufic_action_1"), String("insufic_action_2")},
						MetricName:                       String("my-metric-name"),
						Namespace:                        String("my-namespace"),
						OKActions:                        []*string{String("ok_action_1"), String("ok_action_2")},
						Period:                           Int64(17),
						Statistic:                        String("Average"),
						Threshold:                        Float64(42.17),
						TreatMissingData:                 String("ignore"),
						Unit:                             String("Kilobits"),
					}},
				}, nil
			},
			PutMetricAlarmFunc: func(param0 *cloudwatch.PutMetricAlarmInput) (*cloudwatch.PutMetricAlarmOutput, error) {
				return nil, nil
			},
		}).ExpectInput("PutMetricAlarm", &cloudwatch.PutMetricAlarmInput{
			ActionsEnabled:                   Bool(true),
			AlarmActions:                     []*string{String("old_action_1")},
			AlarmDescription:                 String("my description"),
			AlarmName:                        String("my-alarm-name"),
			ComparisonOperator:               String("GreaterThanOrEqualToThreshold"),
			Dimensions:                       []*cloudwatch.Dimension{{Name: String("Dim1"), Value: String("Val1")}, {Name: String("Dim2"), Value: String("Val2")}},
			EvaluateLowSampleCountPercentile: String("evaluate"),
			EvaluationPeriods:                Int64(42),
			ExtendedStatistic:                String("p0.5"),
			InsufficientDataActions:          []*string{String("insufic_action_1"), String("insufic_action_2")},
			MetricName:                       String("my-metric-name"),
			Namespace:                        String("my-namespace"),
			OKActions:                        []*string{String("ok_action_1"), String("ok_action_2")},
			Period:                           Int64(17),
			Statistic:                        String("Average"),
			Threshold:                        Float64(42.17),
			TreatMissingData:                 String("ignore"),
			Unit:                             String("Kilobits"),
		}).ExpectInput("DescribeAlarms", &cloudwatch.DescribeAlarmsInput{
			AlarmNames: []*string{String("my-alarm-to-detach")},
		}).
			ExpectCalls("PutMetricAlarm", "DescribeAlarms").Run(t)
	})
}
