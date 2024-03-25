package gen

import "os"

const (
	configYamlDefault = `
core:
  ignoreSystemCheck: false # ignore system
  arch: arm64 # amd64 arm64
  registry: sit-registry.qm.cn/qkp-system/kubernetes # image registry address
  networkAdapter: eth0 # network interface name
containerRuntime:
  type: containerd # container runtime type (containerd , docker)
  version: v1.6.21 # container runtime version
  insecureRegistryList: # container runtime insecure registry list
    - sit-registry.qm.cn
etcd:
  version: v3.5.0 # etcd version
  rootPath: /var/lib/etcd # etcd root path
  replicas: 3 # default replicas is master node number , value must in （1，3，5，7）
kubernetes:
  version: v1.21.5 # kubernetes version
  serviceCidr: 10.91.0.0/16  # service cidr address
  podCidr: 10.241.0.0/16 # pod cidr address
kubeVip:
  enable: true # enable kube-vip
  vip: 10.211.55.200 # vip address
  version: v0.3.8 # kube vip version
cni:
  enable: true # enable cni
  type: calico # cni type (calico)
  version: v3.20.0 # cni version
`
	hostYamlFileDefault = `
masters:
  - hostname: m01
    ip: 10.211.55.7
    username: root
    password: xiang1234
    port: 22
  - hostname: m02
    ip: 10.211.55.8
    username: root
    password: xiang1234
    port: 22
  - hostname: m03
    ip: 10.211.55.9
    username: root
    password: xiang1234
    port: 22
workers:
  - hostname: n01
    ip: 10.211.55.10
    username: root
    password: xiang1234
    port: 22
  - hostname: n02
    ip: 10.211.55.11
    username: root
    password: xiang1234
    port: 22
`
)

func GenDefaultConfig() error {
	err := os.WriteFile("host.yaml", []byte(hostYamlFileDefault), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile("config.yaml", []byte(configYamlDefault), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
