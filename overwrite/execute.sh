#!/bin/bash

BASEDIR=$(dirname $0)

cd $BASEDIR
cp ./flex.go ../vendor/github.com/line/line-bot-sdk-go/linebot/flex.go
