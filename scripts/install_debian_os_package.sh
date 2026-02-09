#!/bin/bash

grep -rq Apt::Cmd::Disable-Script-Warning /etc/apt/apt.conf.d/ || \
  echo "Apt::Cmd::Disable-Script-Warning true;" > /etc/apt/apt.conf.d/90disablescriptwarning

export DEBIAN_FRONTEND=noninteractive
apt update -yqq
apt install -yqq --no-install-recommends \
  vim curl wget git netcat-openbsd \
  rsync tcpdump dnsutils net-tools traceroute jq ripgrep \
  build-essential gcc less uuid-runtime
