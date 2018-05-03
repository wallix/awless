package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/elb"
)

func TestClassicLoadbalancer(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Run("without healthcheck", func(t *testing.T) {
			Template("create classicloadbalancer name=my-new-loadbalancer listeners=[HTTPS:443:HTTP:80,TCP:123456:UDP:8080] zones=[us-west-2b,us-west-2c] subnets=[sub-1,sub-2] securitygroups=[sg-1,sg-2] tags=Env:Prod,Dept:Product").Mock(&elbMock{
				CreateLoadBalancerFunc: func(input *elb.CreateLoadBalancerInput) (*elb.CreateLoadBalancerOutput, error) {
					return &elb.CreateLoadBalancerOutput{
						DNSName: String(""),
					}, nil

				},
				ConfigureHealthCheckFunc: func(*elb.ConfigureHealthCheckInput) (*elb.ConfigureHealthCheckOutput, error) {
					return nil, nil // ignored
				}}).
				ExpectInput("CreateLoadBalancer", &elb.CreateLoadBalancerInput{
					AvailabilityZones: []*string{String("us-west-2b"), String("us-west-2c")},
					LoadBalancerName:  String("my-new-loadbalancer"),
					Subnets:           []*string{String("sub-1"), String("sub-2")},
					Tags:              []*elb.Tag{{Key: String("Env"), Value: String("Prod")}, {Key: String("Dept"), Value: String("Product")}},
					Listeners: []*elb.Listener{
						{Protocol: String("HTTPS"), LoadBalancerPort: Int64(443), InstanceProtocol: String("HTTP"), InstancePort: Int64(80)},
						{Protocol: String("TCP"), LoadBalancerPort: Int64(123456), InstanceProtocol: String("UDP"), InstancePort: Int64(8080)},
					},
					SecurityGroups: []*string{String("sg-1"), String("sg-2")},
				}).
				ExpectInput("ConfigureHealthCheck", &elb.ConfigureHealthCheckInput{
					LoadBalancerName: String("my-new-loadbalancer"),
					HealthCheck: &elb.HealthCheck{
						HealthyThreshold:   Int64(10),
						UnhealthyThreshold: Int64(2),
						Interval:           Int64(30),
						Timeout:            Int64(5),
						Target:             String("HTTP:80/"),
					},
				}).ExpectCalls("CreateLoadBalancer", "ConfigureHealthCheck").ExpectCommandResult("my-new-loadbalancer").
				ExpectRevert("delete classicloadbalancer name=my-new-loadbalancer").Run(t)
		})

		t.Run("with healthcheck", func(t *testing.T) {
			Template("create classicloadbalancer healthcheck-path=index.html name=my-new-loadbalancer listeners=HTTPS:443:HTTP:8080 zones=[us-west-2b,us-west-2c] subnets=[sub-1,sub-2] securitygroups=[sg-1,sg-2] tags=Env:Prod,Dept:Product").Mock(&elbMock{
				CreateLoadBalancerFunc: func(*elb.CreateLoadBalancerInput) (*elb.CreateLoadBalancerOutput, error) {
					return &elb.CreateLoadBalancerOutput{
						DNSName: String(""),
					}, nil
				},
				ConfigureHealthCheckFunc: func(*elb.ConfigureHealthCheckInput) (*elb.ConfigureHealthCheckOutput, error) {
					return nil, nil // ignored
				}}).
				ExpectInput("CreateLoadBalancer", &elb.CreateLoadBalancerInput{
					AvailabilityZones: []*string{String("us-west-2b"), String("us-west-2c")},
					LoadBalancerName:  String("my-new-loadbalancer"),
					Subnets:           []*string{String("sub-1"), String("sub-2")},
					Tags:              []*elb.Tag{{Key: String("Env"), Value: String("Prod")}, {Key: String("Dept"), Value: String("Product")}},
					Listeners: []*elb.Listener{
						{Protocol: String("HTTPS"), LoadBalancerPort: Int64(443), InstanceProtocol: String("HTTP"), InstancePort: Int64(8080)},
					},
					SecurityGroups: []*string{String("sg-1"), String("sg-2")},
				}).
				ExpectInput("ConfigureHealthCheck", &elb.ConfigureHealthCheckInput{
					LoadBalancerName: String("my-new-loadbalancer"),
					HealthCheck: &elb.HealthCheck{
						HealthyThreshold:   Int64(10),
						UnhealthyThreshold: Int64(2),
						Interval:           Int64(30),
						Timeout:            Int64(5),
						Target:             String("HTTP:8080/index.html"),
					},
				}).ExpectCalls("CreateLoadBalancer", "ConfigureHealthCheck").ExpectCommandResult("my-new-loadbalancer").
				ExpectRevert("delete classicloadbalancer name=my-new-loadbalancer").Run(t)
		})
	})

	t.Run("update", func(t *testing.T) {
		Template("update classicloadbalancer name=my-classic-loadb health-interval=1 health-timeout=2 healthy-threshold=3 unhealthy-threshold=4 health-target=HTTP:80/home.html").Mock(&elbMock{
			ConfigureHealthCheckFunc: func(*elb.ConfigureHealthCheckInput) (*elb.ConfigureHealthCheckOutput, error) {
				return nil, nil // ignored
			}}).
			ExpectInput("ConfigureHealthCheck", &elb.ConfigureHealthCheckInput{
				LoadBalancerName: String("my-classic-loadb"),
				HealthCheck: &elb.HealthCheck{
					HealthyThreshold:   Int64(3),
					UnhealthyThreshold: Int64(4),
					Interval:           Int64(1),
					Timeout:            Int64(2),
					Target:             String("HTTP:80/home.html"),
				},
			}).ExpectCalls("ConfigureHealthCheck").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete classicloadbalancer name=my-classic-loadb").Mock(&elbMock{
			DeleteLoadBalancerFunc: func(input *elb.DeleteLoadBalancerInput) (*elb.DeleteLoadBalancerOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteLoadBalancer", &elb.DeleteLoadBalancerInput{
				LoadBalancerName: String("my-classic-loadb"),
			}).ExpectCalls("DeleteLoadBalancer").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		Template("attach classicloadbalancer name=my-classic-loadb instance=i-123456").Mock(&elbMock{
			RegisterInstancesWithLoadBalancerFunc: func(input *elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {
				return &elb.RegisterInstancesWithLoadBalancerOutput{
					Instances: []*elb.Instance{{InstanceId: String("i-123456")}},
				}, nil
			}}).
			ExpectInput("RegisterInstancesWithLoadBalancer", &elb.RegisterInstancesWithLoadBalancerInput{
				LoadBalancerName: String("my-classic-loadb"),
				Instances:        []*elb.Instance{{InstanceId: String("i-123456")}},
			}).ExpectCalls("RegisterInstancesWithLoadBalancer").ExpectCommandResult("i-123456").
			ExpectRevert("detach classicloadbalancer instance=i-123456 name=my-classic-loadb").Run(t)
	})

	t.Run("detach", func(t *testing.T) {
		Template("detach classicloadbalancer name=my-classic-loadb instance=i-123456").Mock(&elbMock{
			DeregisterInstancesFromLoadBalancerFunc: func(input *elb.DeregisterInstancesFromLoadBalancerInput) (*elb.DeregisterInstancesFromLoadBalancerOutput, error) {
				return &elb.DeregisterInstancesFromLoadBalancerOutput{
					Instances: []*elb.Instance{{InstanceId: String("i-123456")}},
				}, nil
			}}).
			ExpectInput("DeregisterInstancesFromLoadBalancer", &elb.DeregisterInstancesFromLoadBalancerInput{
				LoadBalancerName: String("my-classic-loadb"),
				Instances:        []*elb.Instance{{InstanceId: String("i-123456")}},
			}).ExpectCalls("DeregisterInstancesFromLoadBalancer").ExpectCommandResult("i-123456").
			ExpectRevert("attach classicloadbalancer instance=i-123456 name=my-classic-loadb").Run(t)
	})
}
