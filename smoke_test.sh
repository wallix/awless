#!/usr/bin/env bash
TMP_FILE=./tmp-integration-test.awless
BIN=./awless

echo "building awless"
go build

echo "flushing awless logs..."
$BIN log --delete

INSTANCE_NAME=awless-integration-test-`date +%s`

/bin/cat > $TMP_FILE <<EOF
testvpc = create vpc cidr=10.0.0.0/24
testsubnet = create subnet cidr=10.0.0.0/25 vpc=\$testvpc
create instance subnet=\$testsubnet image={instance.image} type=t2.nano count={instance.count} name=$INSTANCE_NAME
EOF

cat $TMP_FILE

$BIN run ./$TMP_FILE

REVERT_ID=`$BIN log --porcelain | head -1 | cut -f2`

$BIN revert $REVERT_ID -v

rm $TMP_FILE
rm $BIN