#!/bin/bash

set -ex

PARENT_PATH=$(dirname $(cd $(dirname $0); pwd -P))

pushd $PARENT_PATH
mkdir -p build
go build -o build/reputation-adapter adapter/main.go
popd