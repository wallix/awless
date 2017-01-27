package aws

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

func (d *AwsDriver) Update_Securitygroup_DryRun(params map[string]interface{}) (interface{}, error) {
	ipPerms, err := buildIpPermissionsFromParams(params)
	if err != nil {
		return nil, err
	}
	var input interface{}
	if action, ok := params["inbound"].(string); ok {
		switch action {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupIngressInput{DryRun: aws.Bool(true), IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupIngressInput{DryRun: aws.Bool(true), IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'inbound' parameter expect 'authorize' or 'revoke', got %s", action)
		}
	}
	if action, ok := params["outbound"].(string); ok {
		switch action {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupEgressInput{DryRun: aws.Bool(true), IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupEgressInput{DryRun: aws.Bool(true), IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'outbound' parameter expect 'authorize' or 'revoke', got %s", action)
		}
	}
	if input == nil {
		return nil, fmt.Errorf("expect either 'inbound' or 'outbound' parameter")
	}

	// Required params
	setField(params["id"], input, "GroupId")

	switch ii := input.(type) {
	case *ec2.AuthorizeSecurityGroupIngressInput:
		_, err = d.ec2.AuthorizeSecurityGroupIngress(ii)
	case *ec2.RevokeSecurityGroupIngressInput:
		_, err = d.ec2.RevokeSecurityGroupIngress(ii)
	case *ec2.AuthorizeSecurityGroupEgressInput:
		_, err = d.ec2.AuthorizeSecurityGroupEgress(ii)
	case *ec2.RevokeSecurityGroupEgressInput:
		_, err = d.ec2.RevokeSecurityGroupEgress(ii)
	}
	if awsErr, ok := err.(awserr.Error); ok {
		switch code := awsErr.Code(); {
		case code == "DryRunOperation", strings.HasSuffix(code, "NotFound"):
			d.logger.Println("full dry run: update securitygroup ok")
			return nil, nil
		}
	}

	d.logger.Printf("dry run: update securitygroup error: %s", err)
	return nil, err
}

func (d *AwsDriver) Update_Securitygroup(params map[string]interface{}) (interface{}, error) {
	ipPerms, err := buildIpPermissionsFromParams(params)
	if err != nil {
		return nil, err
	}
	var input interface{}
	if action, ok := params["inbound"].(string); ok {
		switch action {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupIngressInput{IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupIngressInput{IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'inbound' parameter expect 'authorize' or 'revoke', got %s", action)
		}
	}
	if action, ok := params["outbound"].(string); ok {
		switch action {
		case "authorize":
			input = &ec2.AuthorizeSecurityGroupEgressInput{IpPermissions: ipPerms}
		case "revoke":
			input = &ec2.RevokeSecurityGroupEgressInput{IpPermissions: ipPerms}
		default:
			return nil, fmt.Errorf("'outbound' parameter expect 'authorize' or 'revoke', got %s", action)
		}
	}
	if input == nil {
		return nil, fmt.Errorf("expect either 'inbound' or 'outbound' parameter")
	}

	// Required params
	setField(params["id"], input, "GroupId")

	var output interface{}
	switch ii := input.(type) {
	case *ec2.AuthorizeSecurityGroupIngressInput:
		output, err = d.ec2.AuthorizeSecurityGroupIngress(ii)
	case *ec2.RevokeSecurityGroupIngressInput:
		output, err = d.ec2.RevokeSecurityGroupIngress(ii)
	case *ec2.AuthorizeSecurityGroupEgressInput:
		output, err = d.ec2.AuthorizeSecurityGroupEgress(ii)
	case *ec2.RevokeSecurityGroupEgressInput:
		output, err = d.ec2.RevokeSecurityGroupEgress(ii)
	}
	if err != nil {
		d.logger.Printf("update securitygroup error: %s", err)
		return nil, err
	}

	d.logger.Println("update securitygroup done")
	return output, nil
}

func buildIpPermissionsFromParams(params map[string]interface{}) ([]*ec2.IpPermission, error) {
	if _, ok := params["cidr"].(string); !ok {
		return nil, fmt.Errorf("invalid cidr '%v'", params["cidr"])
	}
	ipPerm := &ec2.IpPermission{
		IpRanges: []*ec2.IpRange{{CidrIp: aws.String(params["cidr"].(string))}},
	}
	if _, ok := params["protocol"].(string); !ok {
		return nil, fmt.Errorf("invalid protocol '%v'", params["protocol"])
	}
	p := params["protocol"].(string)
	if strings.Contains("any", p) {
		ipPerm.FromPort = aws.Int64(int64(-1))
		ipPerm.ToPort = aws.Int64(int64(-1))
		ipPerm.IpProtocol = aws.String("-1")
		return []*ec2.IpPermission{ipPerm}, nil
	}
	ipPerm.IpProtocol = aws.String(p)
	switch ports := params["portrange"].(type) {
	case int:
		ipPerm.FromPort = aws.Int64(int64(ports))
		ipPerm.ToPort = aws.Int64(int64(ports))
	case int64:
		ipPerm.FromPort = aws.Int64(ports)
		ipPerm.ToPort = aws.Int64(ports)
	case string:
		switch {
		case strings.Contains(ports, "any"):
			ipPerm.FromPort = aws.Int64(int64(-1))
			ipPerm.ToPort = aws.Int64(int64(-1))
		case strings.Contains(ports, "-"):
			from, err := strconv.ParseInt(strings.SplitN(ports, "-", 2)[0], 10, 64)
			if err != nil {
				return nil, err
			}
			to, err := strconv.ParseInt(strings.SplitN(ports, "-", 2)[1], 10, 64)
			if err != nil {
				return nil, err
			}
			ipPerm.FromPort = aws.Int64(from)
			ipPerm.ToPort = aws.Int64(to)
		default:
			port, err := strconv.ParseInt(ports, 10, 64)
			if err != nil {
				return nil, err
			}
			ipPerm.FromPort = aws.Int64(port)
			ipPerm.ToPort = aws.Int64(port)
		}
	}

	return []*ec2.IpPermission{ipPerm}, nil
}
