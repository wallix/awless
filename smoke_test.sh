#!/usr/bin/env bash
TMP_FILE=./tmp-integration-test.awless
BIN=./awless

echo "building awless"
go build

echo "flushing awless logs..."
$BIN log --delete

REGION="us-west-1"
AMI="ami-165a0876"

echo "Setting region $REGION, ami $AMI"
$BIN config set aws.region $REGION
$BIN config set instance.image $AMI

INSTANCE_NAME=awless-integration-test-`date +%s`

/bin/cat > $TMP_FILE <<EOF
testvpc = create vpc cidr=10.0.0.0/24
testsubnet = create subnet cidr=10.0.0.0/25 vpc=\$testvpc
testinstance = create instance subnet=\$testsubnet image={instance.image} type=t2.nano count={instance.count}
create tag resource=\$testinstance key=Name value=$INSTANCE_NAME
EOF

$BIN -v run ./$TMP_FILE
REVERT_ID=`$BIN log --porcelain | head -1 | cut -f2`

$BIN ls instances

ALIAS="\@$INSTANCE_NAME"
eval "$BIN check instance id=$ALIAS state=running timeout=20"


$BIN -v revert $REVERT_ID

rm $TMP_FILE
rm $BIN