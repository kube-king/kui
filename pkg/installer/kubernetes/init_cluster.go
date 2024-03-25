package kubernetes

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
	"kube-invention/pkg/installer/global"
	"os"
	"strings"
)

type InitCluster struct {
}

func (i *InitCluster) Exec(config *config.Config) (err error) {

	initNode := config.Hosts.Masters[0]

	err = initEnv(config, initNode)
	if err != nil {
		return err
	}

	t := task.New("Init Cluster", initNode)
	t1 := strings.SplitN(config.KubernetsOption.ServiceCIDR, ".", 4)
	t1 = append(t1[:(len(t1)-1)], "10")
	clusterDNS := strings.Join(t1, ".")
	hostListArray := config.Hosts.GetHostArray("master")

	var vip string
	if config.KubeVipOption.Enable {
		vip = config.KubeVipOption.Vip
		hostListArray = append(hostListArray, vip)
	} else {
		vip = config.Hosts.Masters[0].Ip
	}
	etcdHostList, err := config.Hosts.GetEtcdHostList(config.EtcdOption.Replicas)
	if err != nil {
		return err
	}

	etcdHostListArr := make([]string, 0)
	for _, h := range etcdHostList {
		etcdHostListArr = append(etcdHostListArr, h.Ip)
	}

	tasks := make([]task.Module, 0)
	t.SetEnv(map[string]interface{}{
		"HOST_LIST":              hostListArray,
		"ETED_HOST_LIST":         etcdHostListArr,
		"CONTAINER_RUNTIME_TYPE": config.ContainerRuntimeOption.Type,
		"KUBERNETES_VERSION":     config.KubernetsOption.Version,
		"KUBE_VIP_VERSION":       config.KubeVipOption.Version,
		"REGISTRY":               config.Core.Registry,
		"VIP_ADDRESS":            vip,
		"VIP_INTERFACE":          config.Core.NetworkAdapter,
		"SERVICE_CIDR":           config.KubernetsOption.ServiceCIDR,
		"POD_CIDR":               config.KubernetsOption.PodCIDR,
		"CLUSTER_DNS":            clusterDNS,
		"ETCD_ROOT_PATH":         config.EtcdOption.RootPath,
		"VERSION":                strings.TrimLeft(config.KubernetsOption.Version, "v"),
	})

	if config.KubeVipOption.Enable {
		tasks = append(tasks, &task.Template{
			Title:            "Deploy kube vip",
			TemplateFilePath: fmt.Sprintf("%v/kube-vip-%v.yaml.tpl", constant.TemplatePath, config.KubeVipOption.Version),
			RemoteFilePath:   "/etc/kubernetes/manifests/kube-vip.yaml",
			Force:            true,
		})
	}

	joinMasterCommand := fmt.Sprintf("%v/join-master-command", constant.DataPath)

	tasks = append(tasks,
		&task.Template{
			Title:            "Create kubeadm yaml",
			TemplateFilePath: fmt.Sprintf("%v/kubeadm-%v.yaml.tpl", constant.TemplatePath, config.KubernetsOption.Version),
			RemoteFilePath:   "/etc/kubernetes/kubeadm-config.yaml",
			Force:            true,
		}, &task.Command{
			Title:   "Init Kubernetes Cluster",
			CmdList: []string{"kubeadm init --upload-certs --config /etc/kubernetes/kubeadm-config.yaml"},
			Callback: func(module task.Module, status task.Result) task.Result {
				global.Log.Info(fmt.Sprintf("Init Kubernetes Cluster Result: \n%v\n", status.Output))
				return status
			},
		}, &task.File{
			Title: "Creating Directory /root/.kube",
			Type:  "directory",
			Mode:  0755,
			Paths: []string{"/root/.kube"},
		},
		&task.Command{
			Title: "Copy kube config file",
			CmdList: []string{
				"cat /etc/kubernetes/admin.conf > /root/.kube/config",
			},
		},
		&task.Command{
			Title: "Generate Kubeadm Join File",
			CmdList: []string{
				fmt.Sprintf(`echo "$(kubeadm token create --ttl 0 --print-join-command) --control-plane --certificate-key $(kubeadm init phase upload-certs --upload-certs --config  /etc/kubernetes/kubeadm-config.yaml | grep -v upload-certs)"`),
			},
			Callback: func(module task.Module, status task.Result) task.Result {
				err := os.WriteFile(joinMasterCommand, []byte(status.Output), os.ModePerm)
				if err != nil {
					return task.Result{}
				}
				return status
			},
		},
	)

	_, err = t.Run(tasks...)
	joinCommand, err := os.ReadFile(joinMasterCommand)
	if err != nil {
		return err
	}
	cmdFields := strings.Fields(string(joinCommand))
	workerJoinCommand := strings.Join(cmdFields[:7], " ")

	err = os.WriteFile(fmt.Sprintf("%v/join-worker-command", constant.DataPath), []byte(workerJoinCommand), os.ModePerm)
	if err != nil {
		return err
	}

	return err
}
