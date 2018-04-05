package awsfetch

import (
	"context"
	"sync"

	"strings"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/wallix/awless/aws/conv"
	"github.com/wallix/awless/cloud/rdf"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/graph"
)

func forEachBucketParallel(ctx context.Context, cache fetch.Cache, api s3iface.S3API, f func(b *s3.Bucket) error) error {
	var buckets []*s3.Bucket

	if val, e := cache.Get("getBucketsPerRegion", func() (interface{}, error) {
		return getBucketsPerRegion(ctx, api)
	}); e != nil {
		return e
	} else if v, ok := val.([]*s3.Bucket); ok {
		buckets = v
	}

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

func fetchObjectsForBucket(ctx context.Context, api s3iface.S3API, bucket *s3.Bucket, resourcesC chan<- *graph.Resource) error {
	objectc := make(chan []*s3.Object)
	errc := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := api.ListObjectsPages(&s3.ListObjectsInput{Bucket: bucket.Name}, func(page *s3.ListObjectsOutput, lastPage bool) bool {
			objectc <- page.Contents
			return !lastPage
		}); err != nil {
			errc <- err
			return
		}
	}()

	processObjects := func(objs []*s3.Object) {
		for _, output := range objs {
			res, err := awsconv.NewResource(output)
			if err != nil {
				errc <- err
				return
			}
			res.SetProperty("Bucket", awssdk.StringValue(bucket.Name))
			resourcesC <- res
			parent, err := awsconv.InitResource(bucket)
			if err != nil {
				errc <- err
				return
			}
			res.AddRelation(rdf.ChildrenOfRel, parent)
			resourcesC <- parent
		}
	}

	go func() {
		wg.Wait()
		close(objectc)
		close(errc)
	}()

	for {
		select {
		case err := <-errc:
			return err
		case objects, ok := <-objectc:
			if !ok {
				return nil
			}
			processObjects(objects)
		}
	}
}

func getBucketsPerRegion(ctx context.Context, api s3iface.S3API) ([]*s3.Bucket, error) {
	var buckets []*s3.Bucket

	out, err := api.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return buckets, err
	}

	var userBucketName string
	var hasBucketFilter bool
	if id, hasID := getUserFiltersFromContext(ctx)["id"]; hasID {
		userBucketName = id
		hasBucketFilter = true
	} else if buck, hasBucket := getUserFiltersFromContext(ctx)["bucket"]; hasBucket {
		userBucketName = buck
		hasBucketFilter = true
	}

	if hasBucketFilter {
		for _, b := range out.Buckets {
			if strings.Contains(strings.ToLower(*b.Name), strings.ToLower(userBucketName)) {
				buckets = append(buckets, b)
			}
		}
	} else {
		buckets = out.Buckets
	}

	bucketc := make(chan *s3.Bucket)
	errc := make(chan error)

	var wg sync.WaitGroup

	for _, bucket := range buckets {
		wg.Add(1)
		go func(b *s3.Bucket) {
			defer wg.Done()
			loc, err := api.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: b.Name})
			if err != nil {
				errc <- err
				return
			}

			region, _ := ctx.Value("region").(string)
			switch awssdk.StringValue(loc.LocationConstraint) {
			case "":
				if region == "us-east-1" {
					bucketc <- b
				}
			case region:
				bucketc <- b
			}
		}(bucket)
	}
	go func() {
		wg.Wait()
		close(bucketc)
	}()

	var bucketsInRegion []*s3.Bucket
	for {
		select {
		case err := <-errc:
			if err != nil {
				return bucketsInRegion, err
			}
		case b, ok := <-bucketc:
			if !ok {
				return bucketsInRegion, nil
			}
			bucketsInRegion = append(bucketsInRegion, b)
		}
	}
}

func fetchAndExtractGrantsFn(ctx context.Context, api s3iface.S3API, bucketName string) ([]*graph.Grant, error) {
	acls, err := api.GetBucketAcl(&s3.GetBucketAclInput{Bucket: awssdk.String(bucketName)})
	if err != nil {
		return nil, err
	}
	var grants []*graph.Grant
	for _, acl := range acls.Grants {
		displayName := awssdk.StringValue(acl.Grantee.DisplayName)
		granteeType := awssdk.StringValue(acl.Grantee.Type)
		granteeId := awssdk.StringValue(acl.Grantee.ID)

		if awssdk.StringValue(acl.Grantee.EmailAddress) != "" {
			displayName += "<" + awssdk.StringValue(acl.Grantee.EmailAddress) + ">"
		}
		if granteeType == "Group" {
			granteeId += awssdk.StringValue(acl.Grantee.URI)
		}
		grant := &graph.Grant{
			Permission: awssdk.StringValue(acl.Permission),
			Grantee: graph.Grantee{
				GranteeID:          granteeId,
				GranteeType:        granteeType,
				GranteeDisplayName: displayName,
			},
		}
		grants = append(grants, grant)
	}
	return grants, nil
}
