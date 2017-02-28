#!/usr/bin/env bash
TMP_FILE=./tmp-integration-test.awless

echo "flushing awless logs..."
awless log --delete

INSTANCE_NAME=awless-integration-test-`date +%s`

/bin/cat > $TMP_FILE <<EOF
testvpc = create vpc cidr=10.0.0.0/24
testsubnet = create subnet cidr=10.0.0.0/25 vpc=\$testvpc
create instance subnet=\$testsubnet image={instance.image} type={instance.type} count={instance.count} name=$INSTANCE_NAME
EOF

cat $TMP_FILE

awless run ./$TMP_FILE

REVERT_ID=`awless log --porcelain | head -1 | cut -f2`

awless revert $REVERT_ID -v

rm $TMP_FILE