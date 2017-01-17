#!/bin/bash

set -e

echo "building Linux 386"
GOARCH=386 GOOS=linux go build -ldflags "-s -w"
zip awless-linux-386.zip ./awless

rm ./awless

echo "building Linux amd64"
GOARCH=amd64 GOOS=linux go build -ldflags "-s -w"
zip awless-linux-amd64.zip ./awless

rm ./awless

echo "building darwin 386"
GOARCH=386 GOOS=darwin go build -ldflags "-s -w"
zip awless-darwin-386.zip ./awless

rm ./awless

echo "building darwin amd64"
GOARCH=amd64 GOOS=darwin go build -ldflags "-s -w"
zip awless-darwin-amd64.zip ./awless

rm ./awless

echo "building windows 386"
GOARCH=386 GOOS=windows go build -ldflags "-s -w"
zip awless-windows-386.zip ./awless.exe

rm ./awless.exe

echo "building windows amd64"
GOARCH=amd64 GOOS=windows go build -ldflags "-s -w"
zip awless-windows-amd64.zip ./awless.exe

rm ./awless.exe