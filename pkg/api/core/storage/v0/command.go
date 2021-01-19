package v0

import (
	"github.com/vmmgr/node/pkg/api/core/storage"
	"log"
	"os/exec"
	"strconv"
)

func convertImage(d storage.Convert) error {
	//qemu-img convert -f raw -O qcow2 image.img image.qcow2
	out, err := exec.Command("qemu-img", "convert", "-f", d.SrcType, "-O", d.DstType, d.SrcFile, d.DstFile).Output()
	if err != nil {
		return err
	}
	log.Println(string(out))
	return nil
}

func generateImage(fileType, filePath string, fileSize uint) (string, error) {
	//qemu-img create -f qcow2 file.qcow2 100M
	size := strconv.Itoa(int(fileSize)) + "M"
	out, err := exec.Command("qemu-img", "create", "-f", fileType, filePath, size).Output()
	if err != nil {
		return "", err
	}
	log.Println(string(out))
	return string(out), nil
}

func infoImage(filePath string) (string, error) {
	//qemu-img info file.qcow2
	out, err := exec.Command("qemu-img", "info", filePath).Output()
	if err != nil {
		return "", err
	}
	log.Println(string(out))
	return string(out), nil
}

func capacityExpansion(filePath string, size uint) (string, error) {
	out, err := exec.Command("qemu-img", "resize", filePath, strconv.Itoa(int(size))).Output()
	if err != nil {
		return "", err
	}
	log.Println(string(out))
	return string(out), nil
}
