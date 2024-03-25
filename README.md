# Kui
Kubernetes 部署工具 （Kubernetes Install Tools ）

此项目用于部署kubernetes 集群，同事支持国产化麒麟操作系统,支持x86 64和ARM64,包括部署容器运行时 docker、containerd、etcd、kubernetes、kube-vip、calico,目前支持的kubernetes 版本有 
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

### 快速使用
```shell
  # 需要安装依赖
  yum -y install ipset ipvsadm conntrack socat
  # 下载 amd64 二进制文件
  curl -o kui https://github.com/kube-king/kui/releases/download/v0.1/kui-amd64 && chmod +x ./kui
  # 下载 arm64 二进制文件
  curl -o kui https://github.com/kube-king/kui/releases/download/v0.1/kui-arm64 && chmod +x ./kui
  # 生成配置文件模版和主机清单模版
  ./kui gen config  # 需修改 config.yaml 和 host.yaml 
  # 部署一个kubernetes集群
  ./kui init cluster
```

### 功能特性:
<!-- TOC -->
1. [x] 支持 X86_64位和ARM64位架构
2. [x] 支持 单节点master和多节点master集群部署
3. [x] 支持 系统初始化配置
4. [x] 支持 etcd集群和单节点部署
5. [x] 集成 kube-vip 高可用方案 (如需定制化配置，修改templates/kube-vip.yaml.tpl)
6. [x] 支持 docker 、containerd 部署 (如需有定制化安装,修改:templates/install_docker.sh /templates/install_containerd.sh 即可)
7. [x] 网络CNI 支持 calico （目前默认使用ipip模式，如需更改修改templates/calico-v3.20.0.yaml.tpl 即可）
8. [x] 支持在线部署
<!-- TOC -->

### 配置说明:

<!-- TOC -->
#### 配置文件 config/config.yaml
```yaml
core:
  ignoreSystemCheck: false # ignore check system
  arch: arm64 # amd64 arm64
  registry: registry.cn-hangzhou.aliyuncs.com/kube-king # image registry address
  networkAdapter: eth0 # network interface name
containerRuntime:
  type: containerd # container runtime type (containerd , docker) 
  version: v1.6.21 # container runtime version
  insecureRegistryList: # container runtime insecure registry list
    - registry.cn-hangzhou.aliyuncs.com
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
```

#### host节点清单 config/host.yaml

```yaml
masters:
  - hostname: master01
    ip: 1.1.1.1
    username: root
    password: xxxxxx
    port: 22
  - hostname: master02
    ip: 1.1.1.2
    username: root
    password: xxxxxx
    port: 22
  - hostname: master03
    ip: 1.1.1.3
    username: root
    password: xxxxxx
    port: 22
workers:
  - hostname: node01
    ip: 1.1.1.4
    username: root
    password: xxxxxx
    port: 22
  - hostname: node01
    ip: 1.1.1.5
    username: root
    password: xxxxxx
    port: 22
```

### 捐赠 
如果本项目对您有帮助，不妨请作者喝个咖啡。作者邮箱: xiangjiqiang@qq.com

| 微信  | 支付宝                         |
|-----|-----------------------------|
|   ![](docs/images/wechat.jpg)  | ![](docs/images/alipay.jpg) |





