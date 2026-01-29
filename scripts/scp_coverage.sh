#!/bin/bash

# For this script to work, you need to add the following config into your ~/.ssh/config
# Make sure to update IdentityFile line to point where scalezilla project is
#
#Host dev
#  HostName 127.0.0.1
#  User vagrant
#  Port 2222
#  UserKnownHostsFile /dev/null
#  StrictHostKeyChecking no
#  PasswordAuthentication no
#  IdentityFile scalezilla/scalezilla/.vagrant/machines/dev/virtualbox/private_key
#  IdentitiesOnly yes
#  LogLevel FATAL
#  PubkeyAcceptedKeyTypes +ssh-rsa
#  HostKeyAlgorithms +ssh-rsa

COVERAGE=/tmp/coverage.out
scp dev:${COVERAGE} /tmp

if [ -n "$1" ]
then
  go tool cover -html ${COVERAGE}
fi
