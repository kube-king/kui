package config

import (
	"errors"
	"fmt"
	"kube-invention/pkg/client/ssh_client"
	"strings"
)

type Core struct {
	IgnoreSystemCheck bool   `yaml:"ignoreSystemCheck" json:"ignoreSystemCheck"`
	Arch              string `yaml:"arch" json:"arch"`
	Registry          string `yaml:"registry" json:"registry"`
	NetworkAdapter    string `yaml:"networkAdapter" json:"networkAdapter"`
}

const (
	TaskTypeInit = iota
	TaskTypeJoinMaster
	TaskTypeJoinWorker
)

type Config struct {
	TaskType               int                    `json:"-" yaml:"-"`
	Hosts                  Hosts                  `yaml:"-" json:"-"`
	Core                   Core                   `yaml:"core" json:"core"`
	KubernetsOption        KubernetesOption       `yaml:"kubernetes" json:"kubernetes"`
	ContainerRuntimeOption ContainerRuntimeOption `yaml:"containerRuntime" json:"containerRuntime"`
	EtcdOption             EtcdOption             `yaml:"etcd" json:"etcd"`
	KubeVipOption          KubeVipOption          `yaml:"kubeVip" json:"kubeVip"`
	CniOption              CniOption              `yaml:"cni" json:"cni"`
}

type ContainerRuntimeOption struct {
	Type                 string   `yaml:"type" json:"type"`
	InsecureRegistryList []string `yaml:"insecureRegistryList" json:"insecureRegistryList"`
	Version              string   `yaml:"version" json:"version"`
}

type KubernetesOption struct {
	Version     string `yaml:"version" json:"version"`
	ServiceCIDR string `json:"serviceCidr" yaml:"serviceCidr"`
	PodCIDR     string `yaml:"podCidr" json:"podCidr"`
}

type EtcdOption struct {
	Version  string `yaml:"version" json:"version"`
	Replicas int    `yaml:"replicas" json:"replicas"`
	RootPath string `yaml:"rootPath" json:"rootPath"`
}

type KubeVipOption struct {
	Enable  bool   `yaml:"enable" json:"enable"`
	Version string `yaml:"version" json:"version"`
	Vip     string `json:"vip" yaml:"vip"`
}

type CniOption struct {
	Enable  bool   `json:"enable" yaml:"enable"`
	Type    string `json:"type" yaml:"type"`
	Version string `yaml:"version" json:"version"`
}

type Hosts struct {
	Masters []ssh_client.Config `json:"masters,omitempty" yaml:"masters,omitempty"`
	Workers []ssh_client.Config `json:"workers,omitempty" yaml:"workers,omitempty"`
}

type HostLenFunc func() int

func (h *Hosts) GetHostArray(hostType string) (result []string) {
	var host []ssh_client.Config
	switch hostType {
	case "master":
		host = h.Masters
	case "worker":
		host = h.Workers
	}
	result = make([]string, 0)
	for _, h := range host {
		result = append(result, h.Ip)
	}
	return result
}

func (h *Hosts) GetEtcdHostList(replicas int) ([]ssh_client.Config, error) {

	if replicas > len(h.Masters) {
		return nil, errors.New("host number quantity not sufficient")
	}

	result := make([]ssh_client.Config, 0)
	var l int
	if replicas == 0 {
		l = len(h.Masters)
	} else {
		l = replicas
	}

	for i := 0; i < l; i++ {
		result = append(result, h.Masters[i])
	}

	return result, nil
}

func (h *Hosts) GetEtcdInitCluster(hostList []ssh_client.Config) string {
	result := make([]string, 0)
	for _, host := range hostList {
		result = append(result, fmt.Sprintf("etcd-%v=https://%v:2380", host.Ip, host.Ip))
	}
	return strings.Join(result, ",")
}

func (h *Hosts) GetEndpoints(hostList []ssh_client.Config) string {
	result := make([]string, 0)
	for _, host := range hostList {
		result = append(result, fmt.Sprintf("https://%v:2379", host.Ip))
	}
	return strings.Join(result, ",")
}
