# Node
## VM名前
### 命名規則
[GroupID]-[Option]-[VMName]  

|規則|内容|
|---|---|
|GroupID|グループID|
|Option|NetworkIDなど（現時点では0）|
|VMName|VM名|

(例)  
1-0-TestVM  
-> GroupID(1),Option(0),VMName(TestVM)


### サポートOS
* Ubuntu20.04 LTS

### Install
```
sudo apt install qemu-kvm libvirt-daemon-system libvirt-daemon libvirt-dev libvirt-clients bridge-utils libosinfo-bin libguestfs-tools virt-top cloud-image-utils
```

### libvirtの設定
/etc/libvirt/libvirt.conf
```
unix_sock_group = "libvirt"
unix_sock_ro_perms= "0777"
unix_sock_rw_perms= "0770"
```

### PCIパススルーやUSBパススルーをする際は必要
#### IOMMUの有効化
**/etc/grub/default**
```
#Intel
GRUB_CMDLINE_LINUX_DEFAULT="intel_iommu = on"
#AMD
GRUB_CMDLINE_LINUX_DEFAULT="AMD_iommu = on"
```

### Machine一覧表示
```
kvm -M help　(ubuntu)
/usr/libexec/qemu-kvm -machine help (CentOS)
```

### VNCポート、WebSocketポート
WebSocketポートはVNCポート-500により

