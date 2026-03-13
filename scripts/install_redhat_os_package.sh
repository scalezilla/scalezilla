#!/bin/bash

#dnf update -yq
dnf install -yq vim wget git nc \
  bind-utils net-tools traceroute jq gcc \
  uuid-runtime containerd

# set default shell to bash
which chsh || dnf install -yq util-linux-user
grep vagrant /etc/passwd | grep -q /bin/bash || chsh -s /bin/bash vagrant

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

# https://github.com/opencontainers/runc/releases
RUNC_VERSION=1.4.0
wget -qO /usr/local/sbin/runc https://github.com/opencontainers/runc/releases/download/v${RUNC_VERSION}/runc.${ARCH}
chmod +x /usr/local/sbin/runc

mkdir -p /opt/cni/bin
# https://github.com/containernetworking/plugins/releases
CNI_VERSION=1.9.0
wget -P /tmp -q https://github.com/containernetworking/plugins/releases/download/v${CNI_VERSION}/cni-plugins-linux-${ARCH}-v${CNI_VERSION}.tgz
tar Cxzvf /opt/cni/bin /tmp/cni-plugins-linux-${ARCH}-v${CNI_VERSION}.tgz
rm -f /tmp/cni-plugins-linux-${ARCH}-v${CNI_VERSION}.tgz

