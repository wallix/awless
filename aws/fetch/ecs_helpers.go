package awsfetch

import (
	"context"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/wallix/awless/fetch"
)

func getClusterArns(ctx context.Context, cache fetch.Cache, api ecsiface.ECSAPI) ([]string, error) {
	var arns []string
	if clusterName, hasFilter := getUserFiltersFromContext(ctx)["cluster"]; hasFilter {
		out, err := api.DescribeClusters(&ecs.DescribeClustersInput{Clusters: []*string{&clusterName}})
		if err != nil {
			return arns, err
		}
		for _, c := range out.Clusters {
			arns = append(arns, awssdk.StringValue(c.ClusterArn))
		}
	} else {
		if val, cerr := cache.Get("getClustersNames", func() (interface{}, error) {
			err := api.ListClustersPages(&ecs.ListClustersInput{}, func(out *ecs.ListClustersOutput, lastPage bool) (shouldContinue bool) {
				arns = append(arns, awssdk.StringValueSlice(out.ClusterArns)...)
				return out.NextToken != nil
			})
			return arns, err
		}); cerr != nil {
			return arns, cerr
		} else if v, ok := val.([]string); ok {
			arns = v
		}
	}
	return arns, nil
}

func getAllTasks(ctx context.Context, cache fetch.Cache, api ecsiface.ECSAPI) (res []*ecs.Task, err error) {
	clusterArns, cerr := getClusterArns(ctx, cache, api)
	if cerr != nil {
		return res, cerr
	}

	type listTasksOutput struct {
		err     error
		output  *ecs.ListTasksOutput
		cluster *string
	}
	tasksNamesc := make(chan listTasksOutput)
	var wg sync.WaitGroup

	addTaskContainersFunc := func(cl string) func(*ecs.ListTasksOutput, bool) bool {
		return func(out *ecs.ListTasksOutput, lastPage bool) (shouldContinue bool) {
			tasksNamesc <- listTasksOutput{output: out, cluster: awssdk.String(cl)}
			return out.NextToken != nil
		}
	}

	for _, cluster := range clusterArns {
		wg.Add(1)
		go func(cl string) {
			defer wg.Done()
			if er := api.ListTasksPages(&ecs.ListTasksInput{Cluster: &cl, DesiredStatus: awssdk.String("RUNNING")}, addTaskContainersFunc(cl)); er != nil {
				tasksNamesc <- listTasksOutput{err: er}
			}
		}(cluster)

		wg.Add(1)
		go func(cl string) {
			defer wg.Done()
			if er := api.ListTasksPages(&ecs.ListTasksInput{Cluster: &cl, DesiredStatus: awssdk.String("STOPPED")}, addTaskContainersFunc(cl)); er != nil {
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
