#!/bin/bash

COVERAGE=/tmp/coverage.out
cd /home/vagrant/scalezilla
go mod download

go clean -cache && go test -v -race ./... -coverpkg=./... -coverprofile=${COVERAGE}
if [ $? == 0 ]
then
  RUN=true
fi
sed -i '/scalezilla\/scalezilla\/scalezillapb/d' ${COVERAGE}
if [ "${RUN}" == "true" ]
then
  go tool cover -func ${COVERAGE}
fi
