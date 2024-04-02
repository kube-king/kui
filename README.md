# Kui
Kubernetes 部署工具 （Kubernetes Install Tools ）

此项目用于部署kubernetes 集群，同事支持国产化麒麟操作系统,支持x86 64和ARM64,包括部署容器运行时 docker、containerd、etcd、kubernetes、kube-vip、calico,目前支持的kubernetes 版本:
v1.21.5
v1.22.17
v1.23.17
v1.24.0
v1.25.0
v1.26.0
v1.27.0
v1.28.0


#### 操作系统支持
| 操作系统     | 是否测试通过    | 备注                                     |
|----------|----------|----------------------------------------|
| Centos 7 | 是        | centos 8及其以上版本 ，需要调整iptables-legacy 模式 |
| 麒麟V10 SP2 | 是        | 需要卸载podman 和runc                       |
| Rocky Linux 8.6 | 是        | 需要卸载podman 和runc                       |

说明: 其他操作系统可自行测试

版本:

| kubernetes version | etcd version | containerd version | docker version | calico version |
|--------------------|--------------|--------------------|----------------|----------------|
| v1.21.5            | 3.5.0        | v1.7.0             | v20.10.8       | v3.24          |
| v1.22.17           | 3.5.0        | v1.7.0             | v20.10.8       | v3.24          |
| v1.23.17           | 3.5.0        | v1.7.0             | v20.10.8       | v3.24          |
| v1.24.0            | 3.5.3        | v1.7.0             | no support     | v3.24          |
| v1.25.0            | 3.5.4        | v1.7.0             | no support     | v3.24          |
| v1.26.0            | 3.5.6        | v1.7.0             | no support     | v3.25          |
| v1.27.0            | 3.5.7        | v1.7.0             | no support     | v3.26          |
| v1.28.0            | 3.5.9        | v1.7.0             | no support     | v3.27          |


### 安装集群
```shell
  # 需保证集群所有节点已安装以下依赖
  yum -y install ipset ipvsadm conntrack socat
  # 选择cpu架构
  arch = amd64
  # 下载二进制文件
  curl -o kui https://github.com/kube-king/kui/releases/download/v0.2/kui-${arch} && chmod +x ./kui
  # 生成配置文件模版 ( 需根据实际环境修改 config.yaml)
  ./kui gen config --kubernetes-version=v1.25.0 \
                   --container-runtime-type=containerd \
                   --arch=arm64 \
                   --vip=10.211.55.200
  # 生成host主机清单模版
  ./kui gen host -n hosts # 需根据实际主机清单修改 hosts.yaml
  # 部署一个kubernetes集群
  ./kui init --config config.yaml
```
### 批量添加master节点 (注意：需要保留安装好集群之后的 config.yaml 和data 目录中的生成数据)
```shell
  # 生成 add-master-host.yaml 主机清单
  ./kui gen host -n master
  
  # 开始执行添加master节点
  ./kui add master --config config.yaml --hosts add-master-hosts.yaml
```

### 批量添加worker节点 (注意：需要保留安装好集群之后的 config.yaml 和data 目录中的生成数据)
```shell
  # 生成 add-worker-host.yaml 主机清单
  ./kui gen host -n worker
  
  # 开始执行添加worker节点
  ./kui add worker --config config.yaml --hosts add-worker-hosts.yaml
```

### 功能特性:
<!-- TOC -->
1. [x] 支持 X86_64位和ARM64位架构
2. [x] 支持 单节点master和多节点master集群部署
3. [x] 支持 系统初始化配置
4. [x] 支持 etcd集群和单节点部署 （二进制部署）
5. [x] 集成 kube-vip 高可用方案 (如需定制化配置，修改templates/kube-vip.yaml.tpl)
6. [x] 支持 docker 、containerd 部署 (如需有定制化安装,修改:templates/install_docker.sh /templates/install_containerd.sh 即可)
7. [x] 网络CNI 支持 calico （目前默认使用ipip模式，如需更改修改templates/calico-v3.20.0.yaml.tpl 即可）
8. [x] 支持在线部署
9. [x] kubernetes证书已修改为100年
<!-- TOC -->

### 配置说明:

<!-- TOC -->
#### 配置文件 config.yaml
```yaml
core:
  ignoreSystemCheck: true # 是否进行系统检查
  arch: amd64 # cpu架构选择 目前支持 amd64 和arm64 
  registry: registry.cn-hangzhou.aliyuncs.com/kube-king # 镜像仓库地址
  networkAdapter: eth0 # 网卡接口名称
kubernetes:
  version: V1.28.0 # kubernetes 版本
  serviceCidr: 10.91.0.0/16 # service cidr地址
  podCidr: 10.241.0.0/16 # pod cidr地址
containerRuntime:
  type: containerd # kubernetes 1.24开始不支持docker部署，1.24之前的版本支持docker 20.10.8的版本部署
  insecureRegistryList: # 忽略https 的镜像仓库地址
    - registry.cn-hangzhou.aliyuncs.com 
etcd: # etcd 采用外置二进制部署
  replicas: 3 # etcd 副本数
  rootPath: /var/lib/etcd # etcd 数据
kubeVip:
  enable: true # 是否开启kube-vip
  vip: x.x.x.x # vip地址
cni:
  enable: true  # 是否开启步数cni插件
  type: calico # 目前只集成了calico
```

#### host节点清单 hosts.yaml

```yaml
masters: # master 主机清单
  - hostname: master01
    username: root
    password: xxxx
    ip: 1.1.1.1
    port: 22
  - hostname: master02
    username: root
    password: xxxx
    ip: 1.1.1.2
    port: 22
  - hostname: master03
    username: root
    password: xxxx
    ip: 1.1.1.3
    port: 22
workers: # worker 主机清单
  - hostname: worker01
    username: root
    password: xxxx
    ip: 1.1.1.4
    port: 22
  - hostname: worker02
    username: root
    password: xxxx
    ip: 1.1.1.5
    port: 22
```

### 捐赠 
如果本项目对您有帮助，不妨请作者喝个咖啡。作者邮箱: xiangjiqiang@qq.com

| 微信  | 支付宝                         |
|-----|-----------------------------|
|   ![](docs/images/wechat.jpg)  | ![](docs/images/alipay.jpg) |





