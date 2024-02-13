#!/bin/sh
set -e
cd /go/src/proxz
go get .
go build -ldflags="-s -w" -o proxz .
upx proxz
