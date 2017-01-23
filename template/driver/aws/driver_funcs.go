package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (d *AwsDriver) Create_Tags_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateTagsInput{}

	input.DryRun = aws.Bool(true)
	input.Resources = append(input.Resources, aws.String(fmt.Sprint(params["resource"])))

	for k, v := range params {
		if k == "resource" {
			continue
		}
		input.Tags = append(input.Tags, &ec2.Tag{Key: aws.String(k), Value: aws.String(fmt.Sprint(v))})
	}
	_, err := d.api.CreateTags(input)

	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			d.logger.Println("dry run: create tags ok")
			return nil, nil
		}
	}

	d.logger.Printf("dry run: create tags error: %s", err)
	return nil, err
}

func (d *AwsDriver) Create_Tags(params map[string]interface{}) (interface{}, error) {
	input := &ec2.CreateTagsInput{}

	input.Resources = append(input.Resources, aws.String(fmt.Sprint(params["resource"])))

	for k, v := range params {
		if k == "resource" {
			continue
		}
		input.Tags = append(input.Tags, &ec2.Tag{Key: aws.String(k), Value: aws.String(fmt.Sprint(v))})
	}
	_, err := d.api.CreateTags(input)

	if err != nil {
		d.logger.Printf("create tags error: %s", err)
		return nil, err
	}
	d.logger.Println("create tags done")

	return nil, nil
}
