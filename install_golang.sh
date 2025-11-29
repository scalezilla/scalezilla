#!/bin/bash

GO_VERSION=1.25.4
OS=$(uname -s | tr -s "A-Z" "a-z")
OS_ARCH=$(uname -m)
if [[ "${OS}" == "linux" ]]
then
  if [[ "${OS_ARCH}" == "aarch64" ]]
  then
    ARCH=arm64
  else
    ARCH=amd64
  fi
fi
rm -rf /usr/local/go ~/.cache/go-build/* ~/go/pkg/mod/* ~/go/src/* /tmp/go${GO_VERSION}.${OS}-${ARCH}.tar.gz
wget -q https://golang.org/dl/go${GO_VERSION}.${OS}-${ARCH}.tar.gz -P /tmp
tar -C /usr/local -xzf /tmp/go${GO_VERSION}.${OS}-${ARCH}.tar.gz
rm -rf /tmp/go${GO_VERSION}.${OS}-${ARCH}.tar.gz

