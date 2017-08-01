#!/bin/sh

APP_NAME='gomo'
APP_REPO="github.com/tokenshift/gomo"

IFS='/'

set -x

go tool dist list | while read os arch; do
	env GOOS=$os GOARCH=$arch go build -o "${APP_NAME}.${os}_${arch}" "$APP_REPO"
done
