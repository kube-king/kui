package etcd

const (
	CreateEtcdConf = `
cat > {{ .ETCD_PATH }}/cfg/etcd.conf << EOF
#[Member]
ETCD_NAME="etcd-{{ .ip }}"
ETCD_DATA_DIR="{{ .ETCD_PATH }}/data/default.etcd"
ETCD_LISTEN_PEER_URLS="https://{{ .ip }}:2380"
ETCD_LISTEN_CLIENT_URLS="https://{{ .ip }}:2379"

#[Cluster]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://{{ .ip }}:2380"
ETCD_ADVERTISE_CLIENT_URLS="https://{{ .ip }}:2379"
ETCD_INITIAL_CLUSTER="{{ .ETCD_INITIAL_CLUSTER }}"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_INITIAL_CLUSTER_STATE="new"
ETCD_HEARTBEAT_INTERVAL={{ .ETCD_HEARTBEAT_INTERVAL }}
ETCD_ELECTION_TIMEOUT={{ .ETCD_ELECTION_TIMEOUT }}


EOF
`

	CreateEtcdSystemdService = `
cat > /usr/lib/systemd/system/etcd.service << EOF
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
EnvironmentFile={{ .ETCD_PATH }}/cfg/etcd.conf
ExecStart=/usr/bin/etcd \
--cert-file={{ .ETCD_PATH }}/ssl/server.pem \
--key-file={{ .ETCD_PATH }}/ssl/server-key.pem \
--peer-cert-file={{ .ETCD_PATH }}/ssl/server.pem \
--peer-key-file={{ .ETCD_PATH }}/ssl/server-key.pem \
--trusted-ca-file={{ .ETCD_PATH }}/ssl/ca.pem \
--peer-trusted-ca-file={{ .ETCD_PATH }}/ssl/ca.pem \
--logger=zap
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF
systemctl daemon-reload
systemctl enable etcd.service
`
	HealthCheckStatus = `
ETCDCTL_API=3 \
/usr/bin/etcdctl \
--cacert={{ .ETCD_PATH }}/ssl/ca.pem \
--cert={{ .ETCD_PATH }}/ssl/server.pem \
--key={{ .ETCD_PATH }}/ssl/server-key.pem \
--endpoints="{{ .ETCD_ENDPOINTS }}"  member list -w table
`
)
