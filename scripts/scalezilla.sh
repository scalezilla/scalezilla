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

SCALEZILLA_HOME_DIR=/home/vagrant/scalezilla
cd ${SCALEZILLA_HOME_DIR}
go mod download
pgrep go && kill $(pgrep -f "main agent") && sleep 2
rm -rf /var/lib/scalezilla
export SCALEZILLA_LOG_LEVEL=trace
nohup go run main.go agent config -f ${SCALEZILLA_HOME_DIR}/cluster/testdata/config/${CONFIG} &> /tmp/scalezilla.log &

