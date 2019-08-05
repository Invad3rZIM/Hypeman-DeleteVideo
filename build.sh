#!/usr/bin/env bash 
set -xe

# install packages and dependencies
go get github.com/labstack/echo
go get github.com/Sirupsen/logrus

go get github.com/aws/aws-sdk-go/aws
go get github.com/aws/aws-sdk-go/aws/session
go get github.com/aws/aws-sdk-go/service/dynamodb


# build command
go build -o bin/application application.go