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
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

type Security interface {
	stsiface.STSAPI
	GetUserId() (string, error)
	GetAccountId() (string, error)
}

type oncer struct {
	sync.Once
	result interface{}
	err    error
}

type security struct {
	stsiface.STSAPI
}

func NewSecu(sess *session.Session) Security {
	return &security{sts.New(sess)}
}

func (s *security) GetUserId() (string, error) {
	output, err := s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return awssdk.StringValue(output.Arn), nil
}

func (s *security) GetAccountId() (string, error) {
	output, err := s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return awssdk.StringValue(output.Account), nil
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
					g.AddResource(res)
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
				g.AddResource(res)
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
		g.AddResource(res)
		if err != nil {
			return fmt.Errorf("build resource for bucket `%s`: %s", awssdk.StringValue(b.Name), err)
		}
		return nil
	})
	return g, buckets, err
}

func (s *Storage) fetch_all_storageobject_graph() (*graph.Graph, []*s3.Object, error) {
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
		res.Properties["BucketName"] = awssdk.StringValue(bucket.Name)
		g.AddResource(res)
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
			res := graph.InitResource(awssdk.StringValue(url), cloud.Queue)
			res.Properties["Id"] = awssdk.StringValue(url)
			attrs, err := s.GetQueueAttributes(&sqs.GetQueueAttributesInput{AttributeNames: []*string{awssdk.String("All")}, QueueUrl: url})
			if e, ok := err.(awserr.RequestFailure); ok && (e.Code() == sqs.ErrCodeQueueDoesNotExist || e.Code() == sqs.ErrCodeQueueDeletedRecently) {
				return
			}
			if err != nil {
				errc <- err
				return
			}
			for k, v := range attrs.Attributes {
				res.Properties[k] = awssdk.StringValue(v)
			}
			g.AddResource(res)
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
			g.AddResource(res)
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
							g.AddResource(res)
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
