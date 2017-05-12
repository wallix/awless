#!/usr/bin/env bash
set -e

TMP_FILE=./create-basic-infra.awl
TMP_USERDATA_FILE=./tmp-user-data.sh

/bin/cat > $TMP_USERDATA_FILE <<EOF
#!/bin/bash
echo "success" > /tmp/awless-ssh-userdata-success.txt
EOF

BIN=./awless-test

echo "building awless"
go build -o $BIN

echo "flushing awless logs..."
$BIN log --delete-all

ORIG_REGION=$($BIN config get aws.region)
ORIG_IMAGE=$($BIN config get instance.image)

REGION="us-west-1"
AMI="ami-165a0876"

echo "Setting region $REGION, ami $AMI"
$BIN config set aws.region $REGION
$BIN config set instance.image $AMI

SUFFIX=integ-test-$(date +%s)
INSTANCE_NAME=inst-$SUFFIX
VPC_NAME=vpc-$SUFFIX
SUBNET_NAME=subnet-$SUFFIX
KEY_NAME=awless-integ-test-key

/bin/cat > $TMP_FILE <<EOF
vpcname = $VPC_NAME
testvpc = create vpc cidr={vpc-cidr} name=\$vpcname
testsubnet = create subnet cidr={sub-cidr} vpc=\$testvpc name=$SUBNET_NAME
gateway = create internetgateway
attach internetgateway id=\$gateway vpc=\$testvpc
update subnet id=\$testsubnet public=true
rtable = create routetable vpc=\$testvpc
attach routetable id=\$rtable subnet=\$testsubnet
create route cidr=0.0.0.0/0 gateway=\$gateway table=\$rtable
sgroupdesc = "authorize SSH from the Internet"
sgroup = create securitygroup vpc=\$testvpc description=\$sgroupdesc name=ssh-from-internet
update securitygroup id=\$sgroup inbound=authorize protocol=tcp cidr=0.0.0.0/0 portrange=22
testkey = create keypair name=$KEY_NAME
instancecount = {instance.count} # testing var assignement from hole
testinstance = create instance subnet=\$testsubnet image={resolved-image} type=t2.nano count=\$instancecount keypair=\$testkey name=$INSTANCE_NAME userdata=$TMP_USERDATA_FILE securitygroup=\$sgroup
create tag resource=\$testinstance key=Env value=Testing
EOF

RESOLVED_AMI=$($BIN search images debian::jessie --id-only)
$BIN run ./$TMP_FILE vpc-cidr=10.0.0.0/24 sub-cidr=10.0.0.0/25 -e -f resolved-image=$RESOLVED_AMI

ALIAS="\@$INSTANCE_NAME"
eval "$BIN check instance id=$ALIAS state=running timeout=20 -f"

echo "Instance is running. Waiting 20s for system boot"
sleep 20 

SSH_CONNECT=$($BIN ssh $INSTANCE_NAME --print-cli --disable-strict-host-keychecking)
echo "Connecting to instance with $SSH_CONNECT"
RESULT=$($SSH_CONNECT 'cat /tmp/awless-ssh-userdata-success.txt')

if [ "$RESULT" != "success" ]; then
	echo "FAIL to read correct token in remote file after ssh to instance"
	exit -1
fi

echo "Reading token in remote file on instance with success"

REVERT_ID=$($BIN log | grep RevertID | cut -d , -f2 | cut -d : -f2)
$BIN revert $REVERT_ID -e -f

echo "Clean up and reverting back to region '$ORIG_REGION' and ami '$ORIG_IMAGE'"

$BIN config set aws.region $ORIG_REGION
$BIN config set instance.image $ORIG_IMAGE

rm $TMP_FILE $TMP_USERDATA_FILE
rm -f ~/.awless/keys/$KEY_NAME.pem
rm $BIN
