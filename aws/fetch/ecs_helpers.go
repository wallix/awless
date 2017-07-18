package awsfetch

import (
	"context"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/wallix/awless/fetch"
)

func getClustersNames(ctx context.Context, api ecsiface.ECSAPI) (res []*string, err error) {
	err = api.ListClustersPages(&ecs.ListClustersInput{}, func(out *ecs.ListClustersOutput, lastPage bool) (shouldContinue bool) {
		res = append(res, out.ClusterArns...)
		return out.NextToken != nil
	})
	return
}

func getAllTasks(ctx context.Context, cache fetch.Cache, api ecsiface.ECSAPI) (res []*ecs.Task, err error) {
	var clusterArns []*string

	if val, e := cache.Get("getClustersNames", func() (interface{}, error) {
		return getClustersNames(ctx, api)
	}); e != nil {
		return nil, e
	} else if v, ok := val.([]*string); ok {
		clusterArns = v
	}

	type listTasksOutput struct {
		err     error
		output  *ecs.ListTasksOutput
		cluster *string
	}
	tasksNamesc := make(chan listTasksOutput)
	var wg sync.WaitGroup

	addTaskContainersFunc := func(cl *string) func(*ecs.ListTasksOutput, bool) bool {
		return func(out *ecs.ListTasksOutput, lastPage bool) (shouldContinue bool) {
			tasksNamesc <- listTasksOutput{output: out, cluster: cl}
			return out.NextToken != nil
		}
	}

	for _, cluster := range clusterArns {
		wg.Add(1)
		go func(cl *string) {
			defer wg.Done()
			if er := api.ListTasksPages(&ecs.ListTasksInput{Cluster: cl, DesiredStatus: awssdk.String("RUNNING")}, addTaskContainersFunc(cl)); er != nil {
				tasksNamesc <- listTasksOutput{err: er}
			}
		}(cluster)

		wg.Add(1)
		go func(cl *string) {
			defer wg.Done()
			if er := api.ListTasksPages(&ecs.ListTasksInput{Cluster: cl, DesiredStatus: awssdk.String("STOPPED")}, addTaskContainersFunc(cl)); er != nil {
				tasksNamesc <- listTasksOutput{err: er}
			}
		}(cluster)
	}

	type describeTasksOutput struct {
		err    error
		output *ecs.DescribeTasksOutput
	}

	tasksc := make(chan describeTasksOutput)
	var tasksWG sync.WaitGroup

	tasksWG.Add(1)
	go func() {
		defer tasksWG.Done()
		for r := range tasksNamesc {
			if r.err != nil {
				tasksc <- describeTasksOutput{err: r.err}
				return
			}
			if len(r.output.TaskArns) == 0 {
				continue
			}

			tasksWG.Add(1)
			go func(arns []*string, cluster *string) {
				defer tasksWG.Done()
				tasksOut, er := api.DescribeTasks(&ecs.DescribeTasksInput{Cluster: cluster, Tasks: arns})
				tasksc <- describeTasksOutput{err: er, output: tasksOut}
			}(r.output.TaskArns, r.cluster)
		}
	}()

	go func() {
		wg.Wait()
		close(tasksNamesc)
		tasksWG.Wait()
		close(tasksc)
	}()

	for r := range tasksc {
		if err = r.err; err != nil {
			return
		}
		res = append(res, r.output.Tasks...)
	}

	return
}
