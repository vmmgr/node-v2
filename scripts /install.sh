#!/bin/bash
#Install
sudo apt -y update
sudo apt -y upgrade
sudo apt -y install kvm kvm-ipxe qemu-common qemu-kvm qemu-keymaps qemu-kvm-extras qemu-system qemu-user qemu-utils qemu-launcher qemulator qemuctl qtemu qemu-kvm-spice
sudo modprobe kvm_intel
sudo modprobe kvm
sudo apt -y install libvirt0 libvirt-bin libvirt-dev libvirt-doc python-libvirt libvirt-ruby virt-manager virt-viewer  virt-goodies virt-top ubuntu-vm-builder cpu-checker bridge-utils
sudo apt -y install cloud-image-utils