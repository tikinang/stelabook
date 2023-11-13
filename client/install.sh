#!/bin/bash

export GOOS=linux
export GOARCH=arm
export GOARM=7
export CGO_ENABLED=1

go build -o stelabook-client client/main.go
cp stelabook-client "$MEDIA_PATH/stelabook-client"
cp client/stelabook.sh "$MEDIA_PATH/applications/stelabook.app"
