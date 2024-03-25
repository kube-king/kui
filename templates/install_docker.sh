#!/bin/bash
# Set Docker Kernel
cat >/etc/sysctl.d/docker.conf << EOF
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
EOF
sysctl -p /etc/sysctl.d/docker.conf

if [ ! -d /etc/docker ] ; then
mkdir /etc/docker
fi

# Config Docker Daemon File
cat > /etc/docker/daemon.json << EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "insecure-registries": {{.INSECURE_REGISTRY_LIST}},
  "data-root": "/var/lib/docker",
  "log-driver": "json-file",
  "log-level": "warn",
  "log-opts": {
    "max-file": "10",
    "max-size": "1000m"
  }
}
EOF

# Config Docker Service
cat > /usr/lib/systemd/system/docker.service << EOF
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network-online.target firewalld.service
Wants=network-online.target

[Service]
Type=notify
EnvironmentFile=-/etc/sysconfig/docker
ExecStart=/usr/bin/dockerd \$DOCKER_EXTRA_ARGS
ExecReload=/bin/kill -s HUP \$MAINPID
LimitNOFILE=65535
LimitNPROC=infinity
LimitCORE=infinity
TimeoutStartSec=0
Delegate=yes
KillMode=process
Restart=always
StartLimitBurst=3
StartLimitInterval=60s

[Install]
WantedBy=multi-user.target
EOF

# Start Docker Service
systemctl daemon-reload
systemctl enable docker
systemctl restart docker
