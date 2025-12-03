#!/bin/bash

#dnf update -yq
dnf install -yq vim wget git nc \
  bind-utils net-tools traceroute jq gcc

# set default shell to bash
which chsh || dnf install -yq util-linux-user
grep vagrant /etc/passwd | grep -q /bin/bash || chsh -s /bin/bash vagrant

