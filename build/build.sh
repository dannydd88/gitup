#!/bin/bash

SCRIPT_DIR="$(dirname "${BASH_SOURCE:-$0}")"
ROOT_DIR=`python -c 'from __future__ import print_function;import os,sys;print(os.path.realpath(sys.argv[1]))' "$SCRIPT_DIR/.."`

cd "$ROOT_DIR"

VERSION=`git describe --tags`
BUILD=`date -u +%Y%m%d.%H%M%S`
SHA=`git rev-parse --short HEAD`

cd "$ROOT_DIR/cmd"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -x -o "$ROOT_DIR/out/linux-amd64/gitup" \
  -ldflags "-s -w -X main.version=$VERSION -X main.build=$BUILD -X main.sha=$SHA"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
  -x -o "$ROOT_DIR/out/darwin-amd64/gitup" \
  -ldflags "-s -w -X main.version=$VERSION -X main.build=$BUILD -X main.sha=$SHA"

cd "$ROOT_DIR/out/linux-amd64"
tar zcf gitup-linux-amd64.tar.gz gitup
cp gitup-linux-amd64.tar.gz "$ROOT_DIR/out/"

cd "$ROOT_DIR/out/darwin-amd64"
tar zcf gitup-darwin-amd64.tar.gz gitup
cp gitup-darwin-amd64.tar.gz "$ROOT_DIR/out/"

cd -
