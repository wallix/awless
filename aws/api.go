/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/cloud/properties"
	"github.com/wallix/awless/graph"
)

func GetCloudServicesForAPIs(apis ...string) (services []cloud.Service) {
	unique := make(map[string]struct{})
	for _, api := range apis {
		if name, ok := ServicePerAPI[api]; ok {
			if service, exists := cloud.ServiceRegistry[name]; exists {
				if _, done := unique[name]; !done {
					unique[name] = struct{}{}
					services = append(services, service)
				}
			}
		}
	}
	return
}

func GetCloudServicesForTypes(types ...string) (services []cloud.Service) {
	unique := make(map[string]struct{})
	for _, typ := range types {
		if name, ok := ServicePerResourceType[typ]; ok {
			if service, exists := cloud.ServiceRegistry[name]; exists {
				if _, done := unique[name]; !done {
					unique[name] = struct{}{}
					services = append(services, service)
				}
			}
		}
	}
	return
}

func ResourceTypesPerServiceName() map[string][]string {
	out := make(map[string][]string)
	for rT, s := range ServicePerResourceType {
		out[s] = append(out[s], rT)
	}
	return out
}

type oncer struct {
	sync.Once
	result interface{}
	err    error
}

var arnResourceInfoRegex = regexp.MustCompile(`(root)|([\w-.]*)/([\w-./]*)`)

type Identity struct {
	Account, Arn, UserId, ResourceType, ResourcePath, Resource string
}

func (i *Identity) IsRoot() bool {
	return i.Resource == "root"
}

func (i *Identity) IsUserType() bool {
	return i.ResourceType == "user"
}

func (s *Access) GetIdentity() (*Identity, error) {
	resp, err := s.STSAPI.GetCallerIdentity(nil)
	if err != nil {
		return nil, err
	}

	ident := &Identity{
		Account: awssdk.StringValue(resp.Account),
		Arn:     awssdk.StringValue(resp.Arn),
		UserId:  awssdk.StringValue(resp.UserId),
	}

	splits := strings.Split(ident.Arn, ":")
	if l := len(splits); l > 0 {
		ident.ResourcePath = splits[l-1]
		matches := arnResourceInfoRegex.FindStringSubmatch(ident.ResourcePath)
		if len(matches) == 4 {
			if matches[1] == "root" {
				ident.Resource = "root"
				ident.ResourceType = "user"
			} else {
				ident.ResourceType = matches[2]
				ident.Resource = matches[3]
			}
		}
	}

	return ident, nil
}

type UserPolicies struct {
	Username string
	Inlined  []string
	Attached []string
	ByGroup  map[string][]string
}

func (s *Access) GetUserPolicies(username string) (*UserPolicies, error) {
	var wg sync.WaitGroup

	all := &UserPolicies{
		Username: username,
		ByGroup:  make(map[string][]string),
	}

	errc := make(chan error, 4)

	wg.Add(1)
	go func() {
		defer wg.Done()
		policies, err := s.ListUserPolicies(&iam.ListUserPoliciesInput{
			UserName: awssdk.String(username),
		})
		if err != nil {
			errc <- err
			return
		}

		for _, name := range policies.PolicyNames {
			all.Inlined = append(all.Inlined, awssdk.StringValue(name))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		attached, err := s.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
			UserName: awssdk.String(username),
		})
		if err != nil {
			errc <- err
			return
		}

		for _, pol := range attached.AttachedPolicies {
			all.Attached = append(all.Attached, awssdk.StringValue(pol.PolicyName))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		groups, err := s.ListGroupsForUser(&iam.ListGroupsForUserInput{
			UserName: awssdk.String(username),
		})
		if err != nil {
			errc <- err
			return
		}

		type result struct {
			group, policy string
		}
		resultC := make(chan result)
		var wgg sync.WaitGroup
		for _, group := range groups.Groups {
			wgg.Add(1)
			go func(name string) {
				defer wgg.Done()

				output, err := s.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
					GroupName: awssdk.String(name),
				})
				if err != nil {
					errc <- err
					return
				}
				for _, pol := range output.AttachedPolicies {
					resultC <- result{group: name, policy: awssdk.StringValue(pol.PolicyName)}
				}
			}(awssdk.StringValue(group.GroupName))
		}

		go func() {
			wgg.Wait()
			close(resultC)
		}()

		for res := range resultC {
			all.ByGroup[res.group] = append(all.ByGroup[res.group], res.policy)
		}
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for e := range errc {
		return all, e
	}

	return all, nil
}

func (s *Access) fetch_all_user_graph() (*graph.Graph, []*iam.UserDetail, error) {
	g := graph.NewGraph()
	var userDetails []*iam.UserDetail

	var wg sync.WaitGroup
	errc := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var badResErr error
		err := s.GetAccountAuthorizationDetailsPages(&iam.GetAccountAuthorizationDetailsInput{
			Filter: []*string{
				awssdk.String(iam.EntityTypeUser),
			},
		},
			func(out *iam.GetAccountAuthorizationDetailsOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.UserDetailList {
					userDetails = append(userDetails, output)
					var res *graph.Resource
					res, badResErr = newResource(output)
					if badResErr != nil {
						return false
					}
					if badResErr = g.AddResource(res); badResErr != nil {
						return false
					}
				}
				return out.Marker != nil
			})
		if err != nil {
			errc <- err
			return
		}
		if badResErr != nil {
			errc <- err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := s.ListUsersPages(&iam.ListUsersInput{}, func(page *iam.ListUsersOutput, lastPage bool) bool {
			for _, user := range page.Users {
				res, badResErr := newResource(user)
				if badResErr != nil {
					return false
				}
				if badResErr = g.AddResource(res); badResErr != nil {
					return false
				}
			}
			return page.Marker != nil
		})
		if err != nil {
			errc <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return g, userDetails, err
		}
	}

	return g, userDetails, nil
}

func (s *Access) fetch_all_policy_graph() (*graph.Graph, []*iam.Policy, error) {
	g := graph.NewGraph()
	var policies []*iam.Policy

	errc := make(chan error)
	policiesc := make(chan *iam.Policy)

	processPagePolicies := func(page *iam.ListPoliciesOutput) bool {
		for _, p := range page.Policies {
			policiesc <- p
			res, rerr := newResource(p)
			if rerr != nil {
				return false
			}
			if rerr = g.AddResource(res); rerr != nil {
				return false
			}
		}
		return page.Marker != nil
	}

	var wg sync.WaitGroup

	// Return all policies that are only attached
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.ListPoliciesPages(&iam.ListPoliciesInput{OnlyAttached: awssdk.Bool(true)},
			func(out *iam.ListPoliciesOutput, lastPage bool) (shouldContinue bool) {
				return processPagePolicies(out)
			})
		if err != nil {
			errc <- err
		}
	}()

	// Return only self managed policies (local scope)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.ListPoliciesPages(&iam.ListPoliciesInput{Scope: awssdk.String("Local")},
			func(out *iam.ListPoliciesOutput, lastPage bool) (shouldContinue bool) {
				return processPagePolicies(out)
			})
		if err != nil {
			errc <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errc)
		close(policiesc)
	}()

	for {
		select {
		case err := <-errc:
			if err != nil {
				return g, policies, err
			}
		case p, ok := <-policiesc:
			if !ok {
				return g, policies, nil
			}
			policies = append(policies, p)
		}
	}
}

// STORAGE

func (s *Storage) fetch_all_bucket_graph() (*graph.Graph, []*s3.Bucket, error) {
	g := graph.NewGraph()
	var buckets []*s3.Bucket
	bucketM := &sync.Mutex{}

	err := s.foreach_bucket_parallel(func(b *s3.Bucket) error {
		bucketM.Lock()
		buckets = append(buckets, b)
		bucketM.Unlock()
		res, err := newResource(b)
		if err != nil {
			return fmt.Errorf("build resource for bucket `%s`: %s", awssdk.StringValue(b.Name), err)
		}
		if err = g.AddResource(res); err != nil {
			return err
		}
		return nil
	})
	return g, buckets, err
}

func (s *Storage) fetch_all_s3object_graph() (*graph.Graph, []*s3.Object, error) {
	g := graph.NewGraph()
	var cloudResources []*s3.Object

	err := s.foreach_bucket_parallel(func(b *s3.Bucket) error {
		return s.fetchObjectsForBucket(b, g)
	})

	return g, cloudResources, err
}

func (s *Storage) fetchObjectsForBucket(bucket *s3.Bucket, g *graph.Graph) error {
	out, err := s.ListObjects(&s3.ListObjectsInput{Bucket: bucket.Name})
	if err != nil {
		return err
	}

	for _, output := range out.Contents {
		res, err := newResource(output)
		if err != nil {
			return err
		}
		res.Properties["Bucket"] = awssdk.StringValue(bucket.Name)
		if err = g.AddResource(res); err != nil {
			return err
		}
		parent, err := initResource(bucket)
		if err != nil {
			return err
		}
		g.AddParentRelation(parent, res)
	}

	return nil
}

func (s *Storage) getBucketsPerRegion() ([]*s3.Bucket, error) {
	var buckets []*s3.Bucket
	out, err := s.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return buckets, err
	}

	bucketc := make(chan *s3.Bucket)
	errc := make(chan error)

	var wg sync.WaitGroup

	for _, bucket := range out.Buckets {
		wg.Add(1)
		go func(b *s3.Bucket) {
			defer wg.Done()
			loc, err := s.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: b.Name})
			if err != nil {
				errc <- err
				return
			}
			switch awssdk.StringValue(loc.LocationConstraint) {
			case "":
				if s.region == "us-east-1" {
					bucketc <- b
				}
			case s.region:
				bucketc <- b
			}
		}(bucket)
	}
	go func() {
		wg.Wait()
		close(bucketc)
	}()

	for {
		select {
		case err := <-errc:
			if err != nil {
				return buckets, err
			}
		case b, ok := <-bucketc:
			if !ok {
				return buckets, nil
			}
			buckets = append(buckets, b)
		}
	}
}

func (s *Storage) foreach_bucket_parallel(f func(b *s3.Bucket) error) error {
	s.once.Do(func() {
		s.once.result, s.once.err = s.getBucketsPerRegion()
	})
	if s.once.err != nil {
		return s.once.err
	}
	buckets := s.once.result.([]*s3.Bucket)

	errc := make(chan error)
	var wg sync.WaitGroup

	for _, output := range buckets {
		wg.Add(1)
		go func(b *s3.Bucket) {
			defer wg.Done()
			if err := f(b); err != nil {
				errc <- err
			}
		}(output)
	}
	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return err
		}
	}

	return nil
}

// QUEUE

func (s *Queue) fetch_all_queue_graph() (*graph.Graph, []*string, error) {
	g := graph.NewGraph()
	var cloudResources []*string
	out, err := s.ListQueues(&sqs.ListQueuesInput{})
	if err != nil {
		return nil, cloudResources, err
	}
	errc := make(chan error)
	var wg sync.WaitGroup

	for _, output := range out.QueueUrls {
		cloudResources = append(cloudResources, output)
		wg.Add(1)
		go func(url *string) {
			defer wg.Done()
			res := graph.InitResource(cloud.Queue, awssdk.StringValue(url))
			res.Properties[properties.ID] = awssdk.StringValue(url)
			attrs, err := s.GetQueueAttributes(&sqs.GetQueueAttributesInput{AttributeNames: []*string{awssdk.String("All")}, QueueUrl: url})
			if e, ok := err.(awserr.RequestFailure); ok && (e.Code() == sqs.ErrCodeQueueDoesNotExist || e.Code() == sqs.ErrCodeQueueDeletedRecently) {
				return
			}
			if err != nil {
				errc <- err
				return
			}
			for k, v := range attrs.Attributes {
				switch k {
				case "ApproximateNumberOfMessages":
					count, err := strconv.Atoi(awssdk.StringValue(v))
					if err != nil {
						errc <- err
					}
					res.Properties[properties.ApproximateMessageCount] = count
				case "CreatedTimestamp":
					if vv := awssdk.StringValue(v); vv != "" {
						timestamp, err := strconv.ParseInt(vv, 10, 64)
						if err != nil {
							errc <- err
						}
						res.Properties[properties.Created] = time.Unix(int64(timestamp), 0)
					}
				case "LastModifiedTimestamp":
					if vv := awssdk.StringValue(v); vv != "" {
						timestamp, err := strconv.ParseInt(vv, 10, 64)
						if err != nil {
							errc <- err
						}
						res.Properties[properties.Modified] = time.Unix(int64(timestamp), 0)
					}
				case "QueueArn":
					res.Properties[properties.Arn] = awssdk.StringValue(v)
				case "DelaySeconds":
					delay, err := strconv.Atoi(awssdk.StringValue(v))
					if err != nil {
						errc <- err
					}
					res.Properties[properties.Delay] = delay
				}

			}
			if err = g.AddResource(res); err != nil {
				errc <- err
				return
			}
		}(output)

	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
			return g, cloudResources, err
		}
	}

	return g, cloudResources, nil

}

func (s *Infra) fetch_all_listener_graph() (*graph.Graph, []*elbv2.Listener, error) {
	g := graph.NewGraph()
	errc := make(chan error)
	resultc := make(chan *elbv2.Listener)
	var wg sync.WaitGroup
	var cloudResources []*elbv2.Listener

	err := s.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{},
		func(out *elbv2.DescribeLoadBalancersOutput, lastPage bool) (shouldContinue bool) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, lb := range out.LoadBalancers {
					err := s.DescribeListenersPages(&elbv2.DescribeListenersInput{LoadBalancerArn: lb.LoadBalancerArn},
						func(out *elbv2.DescribeListenersOutput, lastPage bool) (shouldContinue bool) {
							for _, listen := range out.Listeners {
								resultc <- listen
							}
							return out.NextMarker != nil
						})
					if err != nil {
						errc <- err
					}
				}
			}()
			return out.NextMarker != nil
		})
	if err != nil {
		return g, cloudResources, err
	}

	go func() {
		wg.Wait()
		close(resultc)
	}()

	for {
		select {
		case err := <-errc:
			if err != nil {
				return g, cloudResources, err
			}
		case listener, ok := <-resultc:
			if !ok {
				return g, cloudResources, nil
			}
			cloudResources = append(cloudResources, listener)
			res, err := newResource(listener)
			if err != nil {
				return g, cloudResources, err
			}
			if err = g.AddResource(res); err != nil {
				return g, cloudResources, err
			}
		}
	}
}

func (s *Dns) fetch_all_record_graph() (*graph.Graph, []*route53.ResourceRecordSet, error) {
	g := graph.NewGraph()
	var cloudResources []*route53.ResourceRecordSet
	zonec := make(chan *route53.HostedZone)
	errc := make(chan error)

	go func() {
		err := s.ListHostedZonesPages(&route53.ListHostedZonesInput{},
			func(out *route53.ListHostedZonesOutput, lastPage bool) (shouldContinue bool) {
				for _, output := range out.HostedZones {
					zonec <- output
				}
				return out.NextMarker != nil
			})
		if err != nil {
			errc <- err
		}
		close(zonec)
	}()

	resultc := make(chan *route53.ResourceRecordSet)

	go func() {
		var wg sync.WaitGroup

		for zone := range zonec {
			wg.Add(1)
			go func(z *route53.HostedZone) {
				defer wg.Done()
				err := s.ListResourceRecordSetsPages(&route53.ListResourceRecordSetsInput{HostedZoneId: z.Id},
					func(out *route53.ListResourceRecordSetsOutput, lastPage bool) (shouldContinue bool) {
						for _, output := range out.ResourceRecordSets {
							resultc <- output
							res, err := newResource(output)
							if err != nil {
								errc <- err
							}
							if err = g.AddResource(res); err != nil {
								errc <- err
							}
							parent, err := initResource(z)
							if err != nil {
								errc <- err
							}
							g.AddParentRelation(parent, res)
						}
						return out.NextRecordName != nil
					})
				if err != nil {
					errc <- err
				}
			}(zone)
		}

		go func() {
			wg.Wait()
			close(resultc)
		}()
	}()

	for {
		select {
		case err := <-errc:
			if err != nil {
				return g, cloudResources, err
			}
		case record, ok := <-resultc:
			if !ok {
				return g, cloudResources, nil
			}
			cloudResources = append(cloudResources, record)
		}
	}
}

func (s *Infra) fetch_all_containercluster_graph() (*graph.Graph, []*ecs.Cluster, error) {
	g := graph.NewGraph()
	var cloudResources []*ecs.Cluster

	var badResErr error
	err := s.ListClustersPages(&ecs.ListClustersInput{}, func(out *ecs.ListClustersOutput, lastPage bool) (shouldContinue bool) {
		var clustersOut *ecs.DescribeClustersOutput

		if clustersOut, badResErr = s.ECSAPI.DescribeClusters(&ecs.DescribeClustersInput{Clusters: out.ClusterArns}); badResErr != nil {
			return false
		}

		for _, cluster := range clustersOut.Clusters {
			cloudResources = append(cloudResources, cluster)
			var res *graph.Resource
			if res, badResErr = newResource(cluster); badResErr != nil {
				return false
			}
			if badResErr = g.AddResource(res); badResErr != nil {
				return false
			}
		}
		return out.NextToken != nil
	})
	if err != nil {
		return g, cloudResources, err
	}

	return g, cloudResources, badResErr
}
