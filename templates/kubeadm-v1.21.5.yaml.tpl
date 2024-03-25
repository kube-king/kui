apiVersion: kubeadm.k8s.io/v1beta2
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: {{ .ip }}
  bindPort: 6443
nodeRegistration:
{{ if eq .CONTAINER_RUNTIME_TYPE "docker" }}
  criSocket: /var/run/dockershim.sock
{{ else }}
  criSocket: /run/containerd/containerd.sock
{{ end }}
  name: {{ .hostname }}
  taints:
    - effect: NoSchedule
      key: node-role.kubernetes.io/master
---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
mode: ipvs
---
apiServer:
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
  certSANs:
{{ range .HOST_LIST }}
    - {{ . }}
{{- end }}
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes-cluster
controllerManager: {}
dns:
  type: CoreDNS
etcd:
  external:
    endpoints:
{{ range .ETED_HOST_LIST }}
      - https://{{ . }}:2379
{{- end }}
    caFile: {{ .ETCD_ROOT_PATH}}/ssl/ca.pem
    certFile: {{ .ETCD_ROOT_PATH}}/ssl/server.pem
    keyFile: {{ .ETCD_ROOT_PATH}}/ssl/server-key.pem
imageRepository: {{ .REGISTRY }}
kind: ClusterConfiguration
kubernetesVersion: {{ .KUBERNETES_VERSION }}
controlPlaneEndpoint: {{ .VIP_ADDRESS }}:6443
networking:
  dnsDomain: cluster.local
  serviceSubnet: {{ .SERVICE_CIDR }}
  podSubnet: {{ .POD_CIDR }}
scheduler: {}
---
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  anonymous:
    enabled: false
  webhook:
    cacheTTL: 0s
    enabled: true
  x509:
    clientCAFile: /etc/kubernetes/pki/ca.crt
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 0s
    cacheUnauthorizedTTL: 0s
cgroupDriver: systemd
clusterDNS:
  - {{ .CLUSTER_DNS }}
clusterDomain: cluster.local
cpuManagerReconcilePeriod: 0s
evictionPressureTransitionPeriod: 0s
fileCheckFrequency: 0s
healthzBindAddress: 127.0.0.1
healthzPort: 10248
httpCheckFrequency: 0s
imageMinimumGCAge: 0s
kind: KubeletConfiguration
logging: {}
nodeStatusReportFrequency: 0s
nodeStatusUpdateFrequency: 0s
rotateCertificates: true
runtimeRequestTimeout: 0s
shutdownGracePeriod: 0s
shutdownGracePeriodCriticalPods: 0s
staticPodPath: /etc/kubernetes/manifests
streamingConnectionIdleTimeout: 0s
syncFrequency: 0s
volumeStatsAggPeriod: 0s