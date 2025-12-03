#!/bin/bash

COVERAGE=/tmp/coverage.out
cd /home/vagrant/scalezilla
go mod download

go test -v -race ./... -coverprofile=${COVERAGE}
sed -i '/scalezilla\/scalezilla\/scalezillapb/d' ${COVERAGE}
go tool cover -func ${COVERAGE}

