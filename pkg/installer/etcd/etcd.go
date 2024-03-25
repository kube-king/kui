package etcd

import (
	"fmt"
	"kube-invention/pkg/client/ssh_client/task"
	"kube-invention/pkg/installer/config"
	"kube-invention/pkg/installer/constant"
	"kube-invention/pkg/installer/global"
	"kube-invention/pkg/utils/certificate"
	"kube-invention/pkg/utils/common"
	"os"
	"time"
)

type Etcd struct {
}

func (e *Etcd) Exec(config *config.Config) (err error) {

	etcdWorkDir := fmt.Sprintf("%v/etcd", constant.DataPath)
	exists, err := common.PathExists(etcdWorkDir)
	if err != nil {
		return err
	}
	if exists {
		err := os.RemoveAll(etcdWorkDir)
		if err != nil {
			return err
		}
	}

	common.MkDirs(0775, constant.CertPath)

	etcdHostList, err := config.Hosts.GetEtcdHostList(config.EtcdOption.Replicas)
	if err != nil {
		return err
	}

	etcdHost := make([]string, 0)
	etcdHost = append(etcdHost, "127.0.0.1")

	for _, h := range etcdHostList {
		etcdHost = append(etcdHost, h.Ip)
	}

	err = certificate.NewEtcdCert(certificate.CertConfig{
		Expire: time.Hour * 24 * 365 * 100,
		SubjectConfig: certificate.Subject{
			Country:    "CN",
			CommonName: "etcd",
		},
		CaCertFilePath:   fmt.Sprintf("%v/ca.pem", constant.CertPath),
		CaKeyFilePath:    fmt.Sprintf("%v/ca-key.pem", constant.CertPath),
		EtcdCertFilePath: fmt.Sprintf("%v/server.pem", constant.CertPath),
		EtcdCsrFilePath:  fmt.Sprintf("%v/etcd.csr", constant.CertPath),
		EtcdKeyFilePath:  fmt.Sprintf("%v/server-key.pem", constant.CertPath),
		DnsList:          etcdHost,
	}).GenerateEtcdCert()

	if err != nil {
		return err
	}

	global.Log.Info("init generator etcd ssl cert success!")

	t := task.New("Deploy Etcd", etcdHostList...)
	t.SetEnv(map[string]interface{}{
		"ETCD_PATH":               config.EtcdOption.RootPath,
		"ETCD_HEARTBEAT_INTERVAL": constant.EtcdHeartbeatInterval,
		"ETCD_ELECTION_TIMEOUT":   constant.EtcdElectionTimeout,
		"ETCD_INITIAL_CLUSTER":    config.Hosts.GetEtcdInitCluster(etcdHostList),
		"ETCD_ENDPOINTS":          config.Hosts.GetEndpoints(etcdHostList),
	})

	_, err = t.Run(&task.Unarchive{
		Title:         "Copy Etcd Binary",
		LocalFilePath: fmt.Sprintf("%v/etcd-%v-linux-%v.tar.gz", constant.BinaryPath, config.EtcdOption.Version, config.Core.Arch),
		RemoteDir:     "/usr/bin/",
		Mode:          0755,
		Force:         true,
	}, &task.File{
		Title: "Create Etcd Dir",
		Type:  "directory",
		Owner: "root",
		Group: "root",
		Paths: []string{
			fmt.Sprintf("%v/bin", config.EtcdOption.RootPath),
			fmt.Sprintf("%v/cfg", config.EtcdOption.RootPath),
			fmt.Sprintf("%v/ssl", config.EtcdOption.RootPath),
			fmt.Sprintf("%v/data", config.EtcdOption.RootPath),
			fmt.Sprintf("%v/wal", config.EtcdOption.RootPath),
		},
	}, &task.Copy{
		Title:          "Copy ca.pem",
		LocalFilePath:  fmt.Sprintf("%v/ca.pem", constant.CertPath),
		RemoteFilePath: fmt.Sprintf("%v/ssl/ca.pem", config.EtcdOption.RootPath),
	}, &task.Copy{
		Title:          "Copy ca-key.pem",
		LocalFilePath:  fmt.Sprintf("%v/ca-key.pem", constant.CertPath),
		RemoteFilePath: fmt.Sprintf("%v/ssl/ca-key.pem", config.EtcdOption.RootPath),
	}, &task.Copy{
		Title:          "Copy server-key.pem",
		LocalFilePath:  fmt.Sprintf("%v/server-key.pem", constant.CertPath),
		RemoteFilePath: fmt.Sprintf("%v/ssl/server-key.pem", config.EtcdOption.RootPath),
	}, &task.Copy{
		Title:          "Copy server.pem",
		LocalFilePath:  fmt.Sprintf("%v/server.pem", constant.CertPath),
		RemoteFilePath: fmt.Sprintf("%v/ssl/server.pem", config.EtcdOption.RootPath),
	}, &task.Command{
		Title: "Create etcd.conf config file",
		CmdList: []string{
			CreateEtcdConf,
		},
	}, &task.Command{
		Title: "Create etcd systemd service file",
		CmdList: []string{
			CreateEtcdSystemdService,
		},
	}, &task.Command{
		Title: "Start etcd service",
		CmdList: []string{
			"systemctl restart etcd",
		},
	}, &task.Command{
		Title:   "Check etcd status",
		CmdList: []string{HealthCheckStatus},
	})

	return err
}
