package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/elbv2"
)

func TestLoadbalancer(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create loadbalancer name=my-new-loadbalancer subnets=sub-1234,sub-2345 subnet-mappings=sub-1234:eipalloc-321, sub-2345:eipalloc-678 "+
			"iptype=ipv4 scheme=Internet-facing securitygroups=sg-1234,sg-2345 type=network").Mock(&elbv2Mock{
			CreateLoadBalancerFunc: func(input *elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
				return &elbv2.CreateLoadBalancerOutput{LoadBalancers: []*elbv2.LoadBalancer{
					{LoadBalancerArn: String("arn:of:new:loadbalancer")},
				}}, nil
			}}).
			ExpectInput("CreateLoadBalancer", &elbv2.CreateLoadBalancerInput{
				Name:           String("my-new-loadbalancer"),
				Subnets:        []*string{String("sub-1234"), String("sub-2345")},
				SubnetMappings: []*elbv2.SubnetMapping{{SubnetId: String("sub-1234"), AllocationId: String("eipalloc-321")}, {SubnetId: String("sub-2345"), AllocationId: String("eipalloc-678")}},
				IpAddressType:  String("ipv4"),
				Scheme:         String("Internet-facing"),
				SecurityGroups: []*string{String("sg-1234"), String("sg-2345")},
				Type:           String("network"),
			}).ExpectCalls("CreateLoadBalancer").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete loadbalancer id=arn:of:loadbalancer:to:delete").Mock(&elbv2Mock{
			DeleteLoadBalancerFunc: func(input *elbv2.DeleteLoadBalancerInput) (*elbv2.DeleteLoadBalancerOutput, error) {
				return nil, nil
			}}).
			ExpectInput("DeleteLoadBalancer", &elbv2.DeleteLoadBalancerInput{
				LoadBalancerArn: String("arn:of:loadbalancer:to:delete"),
			}).ExpectCalls("DeleteLoadBalancer").Run(t)
	})

	t.Run("check", func(t *testing.T) {
		Template("check loadbalancer id=arn:of:loadbalancer:to:check state=active timeout=1").Mock(&elbv2Mock{
			DescribeLoadBalancersFunc: func(input *elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
				return &elbv2.DescribeLoadBalancersOutput{LoadBalancers: []*elbv2.LoadBalancer{
					{LoadBalancerArn: String("arn:of:loadbalancer:to:check"), State: &elbv2.LoadBalancerState{Code: String("active")}},
				}}, nil
			}}).
			ExpectInput("DescribeLoadBalancers", &elbv2.DescribeLoadBalancersInput{
				LoadBalancerArns: []*string{String("arn:of:loadbalancer:to:check")},
			}).ExpectCalls("DescribeLoadBalancers").Run(t)
	})
}
