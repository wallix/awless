#!/usr/bin/env bash
set -e

TMP_FILE=./create-basic-infra.awl
TMP_USERDATA_FILE=./tmp-user-data.sh
SUCCESS_KEYWORD=successfull

/bin/cat > $TMP_USERDATA_FILE <<EOF
#!/bin/bash
echo "{{.Variables.ssh_success_keyword}}" > /tmp/awless-ssh-userdata-success.txt
EOF

BIN=./awless-test

echo "Building latest awless..."
go build -o $BIN

$BIN version

ORIG_REGION=$($BIN config get aws.region)
ORIG_IMAGE=$($BIN config get instance.image)

REGION="us-west-1"
AMI="ami-165a0876"

echo "Setting region $REGION, ami $AMI"
$BIN config set aws.region $REGION
$BIN config set instance.image $AMI

DATE=$(date +%s)
SUFFIX=integ-test-$DATE
INSTANCE_NAME=inst-$SUFFIX

KEY_NAME=awless-integ-test-key
GROUP_NAME=awless-integration-tests

KEY_FILE="$HOME/.awless/keys/$KEY_NAME.pem"

if [ -e  $KEY_FILE ]; then 
	echo "Removing pre existing dummy test key ..."
	rm -f $KEY_FILE
fi

/bin/cat > $TMP_FILE <<EOF
ssh_success_keyword = $SUCCESS_KEYWORD
vpcname = vpc-integ-test-{date}
testvpc = create vpc cidr={vpc-cidr} name=\$vpcname
testsubnet = create subnet cidr={sub-cidr} vpc=\$testvpc name="subnet-integ-test-" + {date}
gateway = create internetgateway
attach internetgateway id=\$gateway vpc=\$testvpc
update subnet id=\$testsubnet public=true
rtable = create routetable vpc=\$testvpc
attach routetable id=\$rtable subnet=\$testsubnet
create route cidr=0.0.0.0/0 gateway=\$gateway table=\$rtable
sgroupdesc = "authorize SSH from the Internet"
sgroup = create securitygroup vpc=\$testvpc description=\$sgroupdesc name=ssh-from-internet
update securitygroup id=\$sgroup inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=22
internetGroupName = "http-from-internet"
sgroupInternet = create securitygroup vpc=\$testvpc description=\$internetGroupName name=\$internetGroupName
update securitygroup id=\$sgroupInternet inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=80
testkey = create keypair name=$KEY_NAME
instancecount = {instance.count} # testing var assignement from hole
instanceSecgroups = [\$sgroup,\$sgroupInternet]
testinstance = create instance subnet=\$testsubnet image={resolved-image} type=t2.nano count=\$instancecount keypair=\$testkey name=inst-integ-test-{date} userdata=$TMP_USERDATA_FILE securitygroup=\$instanceSecgroups
create tag resource=\$testinstance key=Env value=Testing
create policy name=AwlessSmokeTestPolicy resource=* action="ec2:Describe*" effect=Allow
create group name=$GROUP_NAME
attach policy service=lambda access=readonly group=$GROUP_NAME
EOF

RESOLVED_AMI=$($BIN search images debian::jessie --id-only)
$BIN run ./$TMP_FILE vpc-cidr=10.0.0.0/24 sub-cidr=10.0.0.0/25 date=$DATE -e -f resolved-image=$RESOLVED_AMI

ALIAS="\@$INSTANCE_NAME"
eval "$BIN check instance id=$ALIAS state=running timeout=20 -f"

echo "Instance is running. Waiting 20s for system boot"
sleep 20 

SSH_CONNECT=$($BIN ssh $INSTANCE_NAME --print-cli --disable-strict-host-keychecking)
echo "Connecting to instance with $SSH_CONNECT"
RESULT=$($SSH_CONNECT 'cat /tmp/awless-ssh-userdata-success.txt')

if [ "$RESULT" != "$SUCCESS_KEYWORD" ]; then
	echo "FAIL to read correct token in remote file after ssh to instance: got $RESULT, want $SUCCESS_KEYWORD"
	exit -1
fi

echo "Reading keyword $SUCCESS_KEYWORD in remote file on instance with success"

REVERT_ID=$($BIN log -n2 --id-only | head -1)
$BIN revert $REVERT_ID -e -f

echo "Clean up and reverting back to region '$ORIG_REGION' and ami '$ORIG_IMAGE'"

$BIN config set aws.region $ORIG_REGION
$BIN config set instance.image $ORIG_IMAGE

rm $TMP_FILE $TMP_USERDATA_FILE
rm -f ~/.awless/keys/$KEY_NAME.pem
rm $BIN
