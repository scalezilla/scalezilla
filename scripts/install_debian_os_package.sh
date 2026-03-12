#!/bin/bash

grep -rq Apt::Cmd::Disable-Script-Warning /etc/apt/apt.conf.d/ || \
  echo "Apt::Cmd::Disable-Script-Warning true;" > /etc/apt/apt.conf.d/90disablescriptwarning

export DEBIAN_FRONTEND=noninteractive
apt update -yqq
apt install -yqq --no-install-recommends \
  vim curl wget git netcat-openbsd \
  rsync tcpdump dnsutils net-tools traceroute jq ripgrep \
  build-essential gcc less uuid-runtime \
  containerd

systemctl daemon-reload
systemctl enable --now containerd

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

# https://github.com/kubernetes-sigs/cri-tools/releases
CRICTL_VERSION=v1.35.0
wget -P /tmp -q https://github.com/kubernetes-sigs/cri-tools/releases/download/$CRICTL_VERSION/crictl-$CRICTL_VERSION-linux-${ARCH}.tar.gz
sudo tar zxvf /tmp/crictl-$CRICTL_VERSION-linux-${ARCH}.tar.gz -C /usr/local/bin
rm -f /tmp/crictl-$CRICTL_VERSION-linux-${ARCH}.tar.gz
cat > /etc/crictl.yaml <<EOF
runtime-endpoint: unix:///run/containerd/containerd.sock
image-endpoint: unix:///run/containerd/containerd.sock
EOF
