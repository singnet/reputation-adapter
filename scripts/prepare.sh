##!/bin/bash
# Install go ethereum 
#go get -u -x github.com/ethereum/go-ethereum
#cd $GOPATH/src/github.com/ethereum/go-ethereum/
#make
#make devtools
# install executable dependencies

set -ex

PARENT_PATH=$(dirname $(cd $(dirname $0); pwd -P))

pushd $PARENT_PATH
pushd vendor/github.com/golang/protobuf
go install .
popd

