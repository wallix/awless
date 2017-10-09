package awsat

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
)

func TestContainerCluster(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		Template("create containercluster name=my-new-containercluster").
			Mock(&ecsMock{
				CreateClusterFunc: func(param0 *ecs.CreateClusterInput) (*ecs.CreateClusterOutput, error) {
					return &ecs.CreateClusterOutput{
						Cluster: &ecs.Cluster{ClusterArn: String("arn:of:my:new:cluster")},
					}, nil
				},
			}).ExpectInput("CreateCluster", &ecs.CreateClusterInput{
			ClusterName: String("my-new-containercluster"),
		}).
			ExpectCommandResult("arn:of:my:new:cluster").ExpectCalls("CreateCluster").Run(t)
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete containercluster id=arn:of:cluster:to:delete").
			Mock(&ecsMock{
				DeleteClusterFunc: func(param0 *ecs.DeleteClusterInput) (*ecs.DeleteClusterOutput, error) {
					return nil, nil
				},
			}).ExpectInput("DeleteCluster", &ecs.DeleteClusterInput{Cluster: String("arn:of:cluster:to:delete")}).
			ExpectCalls("DeleteCluster").Run(t)
	})

}
