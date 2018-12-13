##!/bin/bash
set -ex

PARENT_PATH=$(dirname $(cd $(dirname $0); pwd -P))

pushd $PARENT_PATH
dep ensure -v

# install executable dependencies
pushd vendor/github.com/golang/protobuf
go install ./protoc-gen-go
popd


popd
