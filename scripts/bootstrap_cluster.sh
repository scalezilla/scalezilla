#!/bin/bash

BOOTSTRAP_TOKEN_FILE=~/.scalezilla_bootstrap_token
uuidgen > ${BOOTSTRAP_TOKEN_FILE}
cd /home/vagrant/scalezilla
go run main.go bootstrap cluster --token $(cat ${BOOTSTRAP_TOKEN_FILE})


