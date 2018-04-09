#!/usr/bin/env bash
set -e

# This test is only working when the region eu-west-2 has been prepared
# Creating resources on the cloud takes ~10 mins
# Reverting it takes > 30 mins (ex: a distribution takes ~15mins!)

BIN=./awless-test

echo "Building latest awless..."
go build -o $BIN
$BIN version

ORIG_REGION=$($BIN config get aws.region)

REGION="eu-west-2"
echo "Setting region $REGION"
$BIN config set aws.region $REGION

# Running mastodonte test
$BIN  run -e smoke_tests/test-all-drivers.aws void.ova-file=$GOPATH/src/github.com/wallix/awless/smoke_tests/test.ova \
cloudformation.policy-file=$GOPATH/src/github.com/wallix/awless/smoke_tests/test-cloudformation.policy \
cloudformation.templatefile=$GOPATH/src/github.com/wallix/awless/smoke_tests/test-cloudformation-sample.template \
lambda.zipfile=$GOPATH/src/github.com/wallix/awless/smoke_tests/test-lambda-function.zip \
random.string=$(env LC_CTYPE=C tr -dc "a-zA-Z0-9" < /dev/urandom | head -c 10) \
principal.account=$($BIN whoami --account-only) principal.user=$($BIN whoami --name-only)

$BIN revert $($BIN log -n1 --id-only)

$BIN config set aws.region $ORIG_REGION