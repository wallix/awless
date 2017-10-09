package awsat

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func TestContainerTask(t *testing.T) {
	t.Run("start", func(t *testing.T) {
		t.Run("service", func(t *testing.T) {
			Template("start containertask name=my-new-service cluster=my-cluster-name desired-count=3 type=service "+
				"role=arn:of:container:role deployment-name=prod loadbalancer.container-name=redis loadbalancer.container-port=6379 "+
				"loadbalancer.targetgroup=arn:of:my:targetgroup").
				Mock(&ecsMock{
					CreateServiceFunc: func(param0 *ecs.CreateServiceInput) (*ecs.CreateServiceOutput, error) {
						return &ecs.CreateServiceOutput{
							Service: &ecs.Service{ServiceArn: String("arn:of:my:new:service")},
						}, nil
					},
				}).ExpectInput("CreateService", &ecs.CreateServiceInput{
				TaskDefinition: String("my-new-service"),
				Cluster:        String("my-cluster-name"),
				DesiredCount:   Int64(3),
				Role:           String("arn:of:container:role"),
				ServiceName:    String("prod"),
				LoadBalancers: []*ecs.LoadBalancer{
					{
						ContainerName:  String("redis"),
						ContainerPort:  Int64(6379),
						TargetGroupArn: String("arn:of:my:targetgroup"),
					},
				},
			}).
				ExpectCommandResult("arn:of:my:new:service").ExpectCalls("CreateService").Run(t)
		})
		t.Run("task", func(t *testing.T) {
			Template("start containertask name=my-new-task cluster=my-cluster-name desired-count=3 type=task").
				Mock(&ecsMock{
					RunTaskFunc: func(param0 *ecs.RunTaskInput) (*ecs.RunTaskOutput, error) {
						return &ecs.RunTaskOutput{
							Tasks: []*ecs.Task{{TaskArn: String("arn:of:new:task")}},
						}, nil
					},
				}).ExpectInput("RunTask", &ecs.RunTaskInput{
				TaskDefinition: String("my-new-task"),
				Cluster:        String("my-cluster-name"),
				Count:          Int64(3),
			}).
				ExpectCommandResult("arn:of:new:task").ExpectCalls("RunTask").Run(t)
		})
	})

	t.Run("stop", func(t *testing.T) {
		t.Run("service", func(t *testing.T) {
			Template("stop containertask cluster=my-cluster-name type=service deployment-name=prod").
				Mock(&ecsMock{
					DeleteServiceFunc: func(param0 *ecs.DeleteServiceInput) (*ecs.DeleteServiceOutput, error) {
						return nil, nil
					},
				}).ExpectInput("DeleteService", &ecs.DeleteServiceInput{
				Cluster: String("my-cluster-name"),
				Service: String("prod"),
			}).ExpectCalls("DeleteService").Run(t)
		})

		t.Run("task", func(t *testing.T) {
			Template("stop containertask cluster=my-cluster-name type=task run-arn=arn:task:to:stop").
				Mock(&ecsMock{
					StopTaskFunc: func(param0 *ecs.StopTaskInput) (*ecs.StopTaskOutput, error) {
						return nil, nil
					},
				}).ExpectInput("StopTask", &ecs.StopTaskInput{
				Cluster: String("my-cluster-name"),
				Task:    String("arn:task:to:stop"),
			}).ExpectCalls("StopTask").Run(t)
		})
	})

	t.Run("update", func(t *testing.T) {
		Template("update containertask name=my-service cluster=my-cluster-name deployment-name=prod desired-count=5").
			Mock(&ecsMock{
				UpdateServiceFunc: func(param0 *ecs.UpdateServiceInput) (*ecs.UpdateServiceOutput, error) {
					return nil, nil
				},
			}).ExpectInput("UpdateService", &ecs.UpdateServiceInput{
			TaskDefinition: String("my-service"),
			Cluster:        String("my-cluster-name"),
			Service:        String("prod"),
			DesiredCount:   Int64(5),
		}).ExpectCalls("UpdateService").Run(t)
	})

	t.Run("attach", func(t *testing.T) {
		t.Run("first container in task", func(t *testing.T) {
			Template("attach containertask name=my-task container-name=redis image=redis/redis memory-hard-limit=128 command='redis --start --fake-param' env=User:Jdoe,DbPasswd:VERYSECRET privileged=true workdir=/home ports=6379,8080:80").
				Mock(&ecsMock{
					DescribeTaskDefinitionFunc: func(param0 *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
						return nil, awserr.New("ClientException", "unable to describe task definition", errors.New("task does not exist"))
					},
					RegisterTaskDefinitionFunc: func(param0 *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
						return &ecs.RegisterTaskDefinitionOutput{TaskDefinition: &ecs.TaskDefinition{TaskDefinitionArn: String("arn:of:my:new:definition")}}, nil
					},
				}).ExpectInput("DescribeTaskDefinition", &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: String("my-task"),
			}).ExpectInput("RegisterTaskDefinition", &ecs.RegisterTaskDefinitionInput{
				Family: String("my-task"),
				ContainerDefinitions: []*ecs.ContainerDefinition{
					{
						Name:             String("redis"),
						Image:            String("redis/redis"),
						Memory:           Int64(128),
						Command:          []*string{String("redis"), String("--start"), String("--fake-param")},
						Environment:      []*ecs.KeyValuePair{{Name: String("User"), Value: String("Jdoe")}, {Name: String("DbPasswd"), Value: String("VERYSECRET")}},
						Privileged:       Bool(true),
						WorkingDirectory: String("/home"),
						PortMappings:     []*ecs.PortMapping{{ContainerPort: Int64(6379)}, {ContainerPort: Int64(80), HostPort: Int64(8080)}},
					},
				},
			}).ExpectCommandResult("arn:of:my:new:definition").ExpectCalls("DescribeTaskDefinition", "RegisterTaskDefinition").Run(t)
		})

		t.Run("more containers in existing task", func(t *testing.T) {
			existingContainer := &ecs.ContainerDefinition{
				Name:             String("redis"),
				Image:            String("redis/redis"),
				Memory:           Int64(128),
				Command:          []*string{String("redis"), String("--start"), String("--fake-param")},
				Environment:      []*ecs.KeyValuePair{{Name: String("User"), Value: String("Jdoe")}, {Name: String("DbPasswd"), Value: String("VERYSECRET")}},
				Privileged:       Bool(true),
				WorkingDirectory: String("/home"),
				PortMappings:     []*ecs.PortMapping{{ContainerPort: Int64(6379)}, {ContainerPort: Int64(80), HostPort: Int64(8080)}},
			}

			Template("attach containertask name=my-task container-name=postgresql image=postgresql memory-hard-limit=64 command=postgresql,--port,3306 ports=3306:3306/tcp").
				Mock(&ecsMock{
					DescribeTaskDefinitionFunc: func(param0 *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
						return &ecs.DescribeTaskDefinitionOutput{
							TaskDefinition: &ecs.TaskDefinition{
								ContainerDefinitions: []*ecs.ContainerDefinition{existingContainer},
								Family:               String("my-task"),
								NetworkMode:          String("bridge"),
							},
						}, nil
					},
					RegisterTaskDefinitionFunc: func(param0 *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
						return &ecs.RegisterTaskDefinitionOutput{TaskDefinition: &ecs.TaskDefinition{TaskDefinitionArn: String("arn:of:my:updated:definition")}}, nil
					},
				}).ExpectInput("DescribeTaskDefinition", &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: String("my-task"),
			}).ExpectInput("RegisterTaskDefinition", &ecs.RegisterTaskDefinitionInput{
				ContainerDefinitions: []*ecs.ContainerDefinition{
					existingContainer,
					{
						Name:         String("postgresql"),
						Image:        String("postgresql"),
						Memory:       Int64(64),
						Command:      []*string{String("postgresql"), String("--port"), String("3306")},
						PortMappings: []*ecs.PortMapping{{ContainerPort: Int64(3306), HostPort: Int64(3306), Protocol: String("tcp")}},
					},
				},
				Family:      String("my-task"),
				NetworkMode: String("bridge"),
			}).ExpectCommandResult("arn:of:my:updated:definition").ExpectCalls("DescribeTaskDefinition", "RegisterTaskDefinition").Run(t)
		})
	})

	t.Run("detach", func(t *testing.T) {
		container1Def := &ecs.ContainerDefinition{
			Name:  String("redis"),
			Image: String("redis/redis"),
		}
		container2Def := &ecs.ContainerDefinition{
			Name: String("posgresql"),
		}
		t.Run("at least 2 containers in task", func(t *testing.T) {
			Template("detach containertask name=my-task container-name=posgresql").
				Mock(&ecsMock{
					DescribeTaskDefinitionFunc: func(param0 *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
						return &ecs.DescribeTaskDefinitionOutput{
							TaskDefinition: &ecs.TaskDefinition{
								ContainerDefinitions: []*ecs.ContainerDefinition{container1Def, container2Def},
								Family:               String("my-task"),
							},
						}, nil
					},
					RegisterTaskDefinitionFunc: func(param0 *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
						return nil, nil
					},
				}).ExpectInput("DescribeTaskDefinition", &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: String("my-task"),
			}).ExpectInput("RegisterTaskDefinition", &ecs.RegisterTaskDefinitionInput{
				ContainerDefinitions: []*ecs.ContainerDefinition{container1Def},
				Family:               String("my-task"),
			}).ExpectCalls("DescribeTaskDefinition", "RegisterTaskDefinition").Run(t)
		})

		t.Run("last container in task", func(t *testing.T) {
			Template("detach containertask name=my-task container-name=redis/redis").
				Mock(&ecsMock{
					DescribeTaskDefinitionFunc: func(param0 *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
						return &ecs.DescribeTaskDefinitionOutput{
							TaskDefinition: &ecs.TaskDefinition{
								ContainerDefinitions: []*ecs.ContainerDefinition{container1Def},
								Family:               String("my-task"),
								TaskDefinitionArn:    String("arn:my-task-to-deregister"),
							},
						}, nil
					},
					DeregisterTaskDefinitionFunc: func(param0 *ecs.DeregisterTaskDefinitionInput) (*ecs.DeregisterTaskDefinitionOutput, error) {
						return nil, nil
					},
				}).ExpectInput("DescribeTaskDefinition", &ecs.DescribeTaskDefinitionInput{
				TaskDefinition: String("my-task"),
			}).ExpectInput("DeregisterTaskDefinition", &ecs.DeregisterTaskDefinitionInput{
				TaskDefinition: String("arn:my-task-to-deregister"),
			}).ExpectCalls("DescribeTaskDefinition", "DeregisterTaskDefinition").Run(t)
		})
	})

	t.Run("delete", func(t *testing.T) {
		Template("delete containertask name=my-task-to-delete all-versions=true").
			Mock(&ecsMock{
				ListTaskDefinitionsFunc: func(param0 *ecs.ListTaskDefinitionsInput) (*ecs.ListTaskDefinitionsOutput, error) {
					return &ecs.ListTaskDefinitionsOutput{TaskDefinitionArns: []*string{String("arn:of:task:to:delete")}}, nil
				},
				DeregisterTaskDefinitionFunc: func(param0 *ecs.DeregisterTaskDefinitionInput) (*ecs.DeregisterTaskDefinitionOutput, error) {
					return nil, nil
				},
			}).ExpectInput("ListTaskDefinitions", &ecs.ListTaskDefinitionsInput{
			FamilyPrefix: String("my-task-to-delete"),
		}).ExpectInput("DeregisterTaskDefinition", &ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: String("arn:of:task:to:delete"),
		}).ExpectCalls("ListTaskDefinitions", "DeregisterTaskDefinition").Run(t)

	})

}
