package gen

import (
	"kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/utils/yaml_util"
	"os"
	"strings"
)

var defaultMasterHosts = []ssh_client.Config{
	{
		Hostname: "master01",
		Ip:       "1.1.1.1",
		Port:     22,
		Username: "root",
		Password: "xxxx",
	},
	{
		Hostname: "master02",
		Ip:       "1.1.1.2",
		Port:     22,
		Username: "root",
		Password: "xxxx",
	},
	{
		Hostname: "master03",
		Ip:       "1.1.1.3",
		Port:     22,
		Username: "root",
		Password: "xxxx",
	},
}

var defaultWorkerHosts = []ssh_client.Config{
	{
		Hostname: "worker01",
		Ip:       "1.1.1.4",
		Port:     22,
		Username: "root",
		Password: "xxxx",
	},
	{
		Hostname: "worker02",
		Ip:       "1.1.1.5",
		Port:     22,
		Username: "root",
		Password: "xxxx",
	},
}

func defaultConfig(kubeVersion string, arch string, containerRuntimeType string, vip string) *config.Config {

	conf := &config.Config{
		Core: config.Core{
			IgnoreSystemCheck: true,
			Arch:              arch,
			Registry:          "registry.cn-hangzhou.aliyuncs.com/kube-king",
			NetworkAdapter:    "eth0",
		},
		ContainerRuntimeOption: config.ContainerRuntimeOption{
			Type: containerRuntimeType,
			InsecureRegistryList: []string{
				"registry.cn-hangzhou.aliyuncs.com",
			},
		},
		EtcdOption: config.EtcdOption{
			RootPath: "/var/lib/etcd",
			Replicas: 3,
		},
		KubernetsOption: config.KubernetesOption{
			ServiceCIDR: "10.91.0.0/16",
			PodCIDR:     "10.241.0.0/16",
			Version:     kubeVersion,
		},
		KubeVipOption: config.KubeVipOption{
			Enable:  true,
			Version: "v0.6.4",
			Vip:     vip,
		},
		CniOption: config.CniOption{
			Enable: true,
			Type:   "calico",
		},
	}

	arr := strings.Split(kubeVersion, ".")
	val := strings.Join(arr[:len(arr)-1], ".")

	switch val {
	case "v1.21":
		conf.CniOption.Version = "v3.24.0"
		conf.EtcdOption.Version = "v3.5.0"
	case "v1.22":
		conf.CniOption.Version = "v3.24.0"
		conf.EtcdOption.Version = "v3.5.0"
	case "v1.23":
		conf.CniOption.Version = "v3.24.0"
		conf.EtcdOption.Version = "v3.5.0"
	case "v1.24":
		conf.CniOption.Version = "v3.24.0"
		conf.EtcdOption.Version = "v3.5.3"
		conf.ContainerRuntimeOption.Type = "containerd"
	case "v1.25":
		conf.CniOption.Version = "v3.24.0"
		conf.EtcdOption.Version = "v3.5.4"
		conf.ContainerRuntimeOption.Type = "containerd"
	case "v1.26":
		conf.CniOption.Version = "v3.25.0"
		conf.EtcdOption.Version = "v3.5.6"
		conf.ContainerRuntimeOption.Type = "containerd"
	case "v1.27":
		conf.EtcdOption.Version = "v3.5.7"
		conf.CniOption.Version = "v3.26.0"
		conf.ContainerRuntimeOption.Type = "containerd"
	case "v1.28":
		conf.CniOption.Version = "v3.27.0"
		conf.EtcdOption.Version = "v3.5.9"
		conf.ContainerRuntimeOption.Type = "containerd"
	}

	switch conf.ContainerRuntimeOption.Type {
	case "containerd":
		switch val {
		case "v1.21", "v1.22", "v1.23":
			conf.ContainerRuntimeOption.Version = "v1.6.21"
		case "v1.24", "v1.25", "v1.26", "v1.27", "v1.28":
			conf.ContainerRuntimeOption.Version = "v1.7.0"
		}
	case "docker":
		conf.ContainerRuntimeOption.Version = "v20.10.8"
	}

	return conf
}

func GenDefaultConfig(kubeVersion string, arch string, containerRuntimeType string, vip string) error {
	conf := defaultConfig(kubeVersion, arch, containerRuntimeType, vip)
	err := yaml_util.YamlToFile(conf, "config.yaml", os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func GenHostConfig(name string) (err error) {
	var val string
	conf := config.Hosts{}
	switch name {
	case "hosts":
		val = "hosts.yaml"
		conf.Masters = defaultMasterHosts
		conf.Workers = defaultWorkerHosts
	case "master":
		val = "add-master-hosts.yaml"
		conf.Masters = defaultMasterHosts
	case "worker":
		val = "add-worker-hosts.yaml"
		conf.Workers = defaultWorkerHosts
	}
	err = yaml_util.YamlToFile(conf, val, os.ModePerm)

	return
}
