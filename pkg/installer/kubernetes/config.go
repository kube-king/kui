package kubernetes

import (
	"fmt"
	"kube-invention/pkg/utils/common"
	"log"
)

const (
	KubeletService = `
cat > /usr/lib/systemd/system/kubelet.service << EOF
[Unit]
Description=kubelet: The Kubernetes Node Agent
Documentation=https://kubernetes.io/docs/
Wants=network-online.target
After=network-online.target

[Service]
ExecStart=/usr/bin/kubelet
Restart=always
StartLimitInterval=0
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
`

	KubeadmConf = `
cat > /usr/lib/systemd/system/kubelet.service.d/10-kubeadm.conf << EOF
# Note: This dropin only works with kubeadm and kubelet v1.11+
[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
# This is a file that "kubeadm init" and "kubeadm join" generates at docker, populating the KUBELET_KUBEADM_ARGS variable dynamically
EnvironmentFile=-/var/lib/kubelet/kubeadm-flags.env
# This is a file that the user can use for overrides of the kubelet args as a last resort. Preferably, the user should use
# the .NodeRegistration.KubeletExtraArgs object in the configuration files instead. KUBELET_EXTRA_ARGS should be sourced from this file.
EnvironmentFile=-/etc/sysconfig/kubelet
ExecStart=
ExecStart=/usr/bin/kubelet \$KUBELET_KUBECONFIG_ARGS \$KUBELET_CONFIG_ARGS \$KUBELET_KUBEADM_ARGS \$KUBELET_EXTRA_ARGS
EOF
systemctl daemon-reload
systemctl enable --now kubelet
`
	ClearSingleEtcdCommand = `
ETCDCTL_API=3 /usr/bin/etcdctl --cacert=/data/etcd/ssl/ca.pem \
      --cert=/data/etcd/ssl/server.pem \
      --key=/data/etcd/ssl/server-key.pem \
      --endpoints="https://{{ .MASTER01_IP }}:2379" \
      del / --prefix
`

	ClearEtcdCommand = `
ETCDCTL_API=3 /usr/bin/etcdctl --cacert=/data/etcd/ssl/ca.pem \
      --cert=/data/etcd/ssl/server.pem \
      --key=/data/etcd/ssl/server-key.pem \
      --endpoints="https://{{ .MASTER01_IP }}:2379,https://{{ .MASTER02_IP }}:2379,https://{{ .MASTER03_IP }}:2379" \
      del / --prefix
`
	CheckKubernetesClusterStatus = `
kubectl get pod -n kube-system|grep %v|awk '{print $1}'|xargs  -n 1  kubectl get pod  -n kube-system  -o=jsonpath='{.metadata.name}={.status.phase}|'
`
)

var (
	KubernetesComponentList = []string{
		"kube-apiserver",
		"kube-controller-manager",
		"kube-scheduler",
		"kube-vip",
	}
)

var checkClusterStatus = func(output string) bool {

	log.Println(output)
	isSuccess := false
	resList := common.ExtractStringReg(output, `(?P<type>.*),(?P<status>.*)`)
	if len(resList) <= 0 {
		return isSuccess
	}
	for _, res := range resList {
		log.Println(fmt.Sprintf("%v:%v", res["type"], res["status"]))
		if res["status"] == "Unhealthy" {
			isSuccess = false
			break
		} else {
			isSuccess = true
		}
	}
	return isSuccess
}
