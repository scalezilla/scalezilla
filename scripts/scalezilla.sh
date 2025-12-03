#!/bin/bash

if [ -n "$1" ]
then
  export CONFIG="$1"
else
  usage CONFIG
  echo "CONFIG argument is mantatory"
  echo "Example: $0 config_success_server11.hcl"
  exit 1
fi
shift

cd /home/vagrant/scalezilla
go mod download
pgrep go && kill $(pgrep go) && sleep 2
nohup go run main.go agent config -f cluster/testdata/config/${CONFIG} &> /tmp/scalezilla.log &

