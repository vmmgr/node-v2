#!/bin/bash
#Install
sudo apt -y update
sudo apt -y upgrade
sudo apt -y install qemu-kvm qemu-keymaps qemu-system qemu-user qemu-utils
sudo modprobe kvm_intel
sudo modprobe kvm
sudo apt -y install libvirt0 libvirt-dev libvirt-doc virt-manager virt-viewer virt-top cpu-checker bridge-utils
sudo apt -y install cloud-image-utils