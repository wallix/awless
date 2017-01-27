package aws

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/shell"
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
	_, err := d.ec2.CreateTags(input)

	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			d.logger.Println("full dry run: create tags ok")
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
	_, err := d.ec2.CreateTags(input)

	if err != nil {
		d.logger.Printf("create tags error: %s", err)
		return nil, err
	}
	d.logger.Println("create tags done")

	return nil, nil
}

func (d *AwsDriver) Create_Keypair_DryRun(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ImportKeyPairInput{}

	input.DryRun = aws.Bool(true)
	setField(params["name"], input, "KeyName")

	if params["name"] == "" {
		err := fmt.Errorf("empty 'name' parameter")
		d.logger.Printf("dry run: saving private key error: %s", err)
		return nil, err
	}

	privKeyPath := filepath.Join(config.KeysDir, fmt.Sprint(params["name"])+".pem")
	_, err := os.Stat(privKeyPath)
	if err == nil {
		fileExist := fmt.Errorf("file already exists at path: %s", privKeyPath)
		d.logger.Printf("dry run: saving private key error: %s", fileExist)
		return nil, fileExist
	}

	return nil, nil
}

func (d *AwsDriver) Create_Keypair(params map[string]interface{}) (interface{}, error) {
	input := &ec2.ImportKeyPairInput{}
	setField(params["name"], input, "KeyName")

	d.logger.Printf("Generating locally a RSA 4096 bits keypair...")
	pub, priv, err := shell.GenerateSSHKeyPair(4096)
	if err != nil {
		d.logger.Printf("generating keypair error: %s", err)
		return nil, err
	}
	privKeyPath := filepath.Join(config.KeysDir, fmt.Sprint(params["name"])+".pem")
	_, err = os.Stat(privKeyPath)
	if err == nil {
		fileExist := fmt.Errorf("file already exists at path: %s", privKeyPath)
		d.logger.Printf("saving private key error: %s", fileExist)
		return nil, fileExist
	}
	err = ioutil.WriteFile(privKeyPath, priv, 0400)
	if err != nil {
		d.logger.Printf("saving private key error: %s", err)
		return nil, err
	}
	fmt.Printf("4096 RSA keypair generated locally and stored in '%s'\n", privKeyPath)
	input.PublicKeyMaterial = pub

	output, err := d.ec2.ImportKeyPair(input)
	if err != nil {
		d.logger.Printf("create keypair error: %s", err)
		return nil, err
	}
	id := aws.StringValue(output.KeyName)
	d.logger.Printf("create keypair '%s' done", id)
	return aws.StringValue(output.KeyName), nil
}
