#!/bin/bash

# Set System Limit
cat > /etc/security/limits.conf << EOF
* soft nofile 102400
* hard nofile 102400
* soft nproc 102400
* hard nproc 102400
EOF

# Set Calico for NetworkManager
cat > /etc/NetworkManager/conf.d/calico.conf << EOF
[keyfile]
unmanaged-devices=interface-name:cali*;interface-name:tunl*;interface-name:vxlan.calico;interface-name:vxlan-v6.calico;interface-name:wireguard.cali;interface-name:wg-v6.cali
EOF



# Disabled Selinux
setenforce 0
sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config

# Disabled Swap
swapoff -a
sed -ir 's/.*swap/#&/g' /etc/fstab

# Load Kernel Module
modprobe ip_conntrack
modprobe br_netfilter
modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
modprobe -- nf_conntrack

if [ ! -d /etc/sysconfig/modules/ ] ; then
  mkdir /etc/sysconfig/modules/
fi

cat > /etc/sysconfig/modules/ipvs.modules <<EOF
#!/bin/bash
ipvs_modules="ip_vs ip_vs_lc ip_vs_wlc ip_vs_rr ip_vs_wrr   ip_vs_lblc ip_vs_lblcr ip_vs_dh ip_vs_sh   ip_vs_nq ip_vs_sed ip_vs_ftp nf_conntrack ip_conntrack br_netfilter"
for kernel_module in \${ipvs_modules}; do
    /sbin/modinfo -F filename \${kernel_module} > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        /sbin/modprobe \${kernel_module}
    fi
done
EOF
chmod 755 /etc/sysconfig/modules/ipvs.modules && bash /etc/sysconfig/modules/ipvs.modules && lsmod | grep ip_vs
cat <<EOF | sudo tee /etc/modules-load.d/kubernetes.conf
overlay
br_netfilter
ip_conntrack
EOF