#!/bin/bash
set -ve

if command -v yum; then
	yum -y install git
else
	apt-get -y install git
fi

GOLANG_TAR=go1.8.1.linux-amd64.tar.gz

curl --fail -o $GOLANG_TAR -L https://storage.googleapis.com/golang/$GOLANG_TAR

tar -C /usr/local -xzf $GOLANG_TAR

GOPATH=/go
mkdir -p $GOPATH/src $GOPATH/bin
chmod -R 777 $GOPATH
echo "export GOPATH=$GOPATH" >> /etc/profile
echo "export PATH=/usr/local/go/bin:/go/bin:$PATH" >> /etc/profile
source /etc/profile

go get github.com/wallix/awless