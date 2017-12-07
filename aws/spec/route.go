/* Copyright 2017 WALLIX

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

package awsspec

import (
	"net"

	"github.com/wallix/awless/cloud/graph"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/wallix/awless/logger"
)

type CreateRoute struct {
	_       string `action:"create" entity:"route" awsAPI:"ec2" awsCall:"CreateRoute" awsInput:"ec2.CreateRouteInput" awsOutput:"ec2.CreateRouteOutput" awsDryRun:""`
	logger  *logger.Logger
	graph   cloudgraph.GraphAPI
	api     ec2iface.EC2API
	Table   *string `awsName:"RouteTableId" awsType:"awsstr" templateName:"table" required:""`
	CIDR    *string `awsName:"DestinationCidrBlock" awsType:"awsstr" templateName:"cidr" required:""`
	Gateway *string `awsName:"GatewayId" awsType:"awsstr" templateName:"gateway" required:""`
}

func (cmd *CreateRoute) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *CreateRoute) Validate_CIDR() error {
	_, _, err := net.ParseCIDR(StringValue(cmd.CIDR))
	return err
}

type DeleteRoute struct {
	_      string `action:"delete" entity:"route" awsAPI:"ec2" awsCall:"DeleteRoute" awsInput:"ec2.DeleteRouteInput" awsOutput:"ec2.DeleteRouteOutput" awsDryRun:""`
	logger *logger.Logger
	graph  cloudgraph.GraphAPI
	api    ec2iface.EC2API
	Table  *string `awsName:"RouteTableId" awsType:"awsstr" templateName:"table" required:""`
	CIDR   *string `awsName:"DestinationCidrBlock" awsType:"awsstr" templateName:"cidr" required:""`
}

func (cmd *DeleteRoute) ValidateParams(params []string) ([]string, error) {
	return validateParams(cmd, params)
}

func (cmd *DeleteRoute) Validate_CIDR() error {
	_, _, err := net.ParseCIDR(StringValue(cmd.CIDR))
	return err
}
