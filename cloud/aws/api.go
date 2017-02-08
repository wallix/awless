package aws

import (
	"regexp"
	"sort"
	"sync"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/wallix/awless/graph"
)

var DefaultAMIUsers = []string{"ec2-user", "ubuntu", "centos", "bitnami", "admin", "root"}

func AllRegions() []string {
	var regions sort.StringSlice
	partitions := endpoints.DefaultResolver().(endpoints.EnumPartitions).Partitions()
	for _, p := range partitions {
		for id := range p.Regions() {
			regions = append(regions, id)
		}
	}
	sort.Sort(regions)
	return regions
}

func IsValidRegion(given string) bool {
	reg, _ := regexp.Compile("^(us|eu|ap|sa|ca)\\-\\w+\\-\\d+$")
	regChina, _ := regexp.Compile("^cn\\-\\w+\\-\\d+$")
	regUsGov, _ := regexp.Compile("^us\\-gov\\-\\w+\\-\\d+$")

	return reg.MatchString(given) || regChina.MatchString(given) || regUsGov.MatchString(given)
}

type Security interface {
	stsiface.STSAPI
	GetUserId() (string, error)
	GetAccountId() (string, error)
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

		out, err := s.GetAccountAuthorizationDetails(&iam.GetAccountAuthorizationDetailsInput{
			Filter: []*string{
				awssdk.String(iam.EntityTypeUser),
			},
		})
		if err != nil {
			errc <- err
			return
		}

		for _, output := range out.UserDetailList {
			userDetails = append(userDetails, output)
			res, err := newResource(output)
			if err != nil {
				errc <- err
				return
			}
			g.AddResource(res)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		out, err := s.ListUsers(&iam.ListUsersInput{})
		if err != nil {
			errc <- err
			return
		}

		for _, output := range out.Users {
			res, err := newResource(output)
			if err != nil {
				errc <- err
				return
			}
			g.AddResource(res)
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

func (s *Storage) fetch_all_storageobject_graph() (*graph.Graph, []*s3.Object, error) {
	g := graph.NewGraph()
	var cloudResources []*s3.Object
	var wg sync.WaitGroup
	errc := make(chan error)

	s.foreach_bucket(func(b *s3.Bucket) error {
		wg.Add(1)
		go func(bucket *s3.Bucket) {
			defer wg.Done()
			errc <- s.fetchObjectsForBucket(bucket, g)
		}(b)
		return nil
	})

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

func (s *Storage) foreach_bucket(f func(b *s3.Bucket) error) error {
	out, err := s.fetch_all_bucket()
	if err != nil {
		return err
	}

	for _, output := range out.(*s3.ListBucketsOutput).Buckets {
		err := f(output)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) fetch_all_bucket() (interface{}, error) {
	return s.ListBuckets(&s3.ListBucketsInput{})
}

func (s *Storage) fetch_all_bucket_graph() (*graph.Graph, []*s3.Bucket, error) {
	g := graph.NewGraph()
	out, err := s.fetch_all_bucket()
	if err != nil {
		return nil, nil, err
	}
	var buckets []*s3.Bucket

	errc := make(chan error)
	var wg sync.WaitGroup

	for _, output := range out.(*s3.ListBucketsOutput).Buckets {
		buckets = append(buckets, output)
		wg.Add(1)
		go func(b *s3.Bucket) {
			defer wg.Done()
			res, err := newResource(b)
			if err != nil {
				errc <- err
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
			return g, nil, nil
		}
	}

	return g, buckets, nil
}
