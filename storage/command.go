package storage

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func RunStorageCmd(cmd []string) result {
	out, err := exec.Command("qemu-img", cmd...).Output()
	if err != nil {
		log.Println(err)
		return result{Err: err}
	}

	log.Println(string(out))
	return result{Info: string(out), Err: nil}
}

func createStorageCmd(s storage) result {
	log.Println("----storage create----")
	if s.size < 0 {
		return result{Err: fmt.Errorf("Wrong storage size !! ")}
	}
	//qemu-img create [-f format] filename [size]
	return RunStorageCmd([]string{"create", "-f", s.format, s.path, strconv.Itoa(s.size) + "M"})
}

func deleteStorageCmd(s storage) result {
	if fileExistsCheck(s.path) {
		if err := os.Remove(s.path); err != nil {
			log.Println(err)
			return result{Info: "Error: failed file delete", Err: err}
		}
	}
	return result{Info: "OK", Err: nil}
}

func ResizeStorageCmd(s storage) result {
	//qemu-img resize [filename] [size]
	return RunStorageCmd([]string{"qemu-img", "resize", s.path, strconv.Itoa(s.size) + "M"})
}

func renameStorageCmd(src, dst string) result {
	if err := os.Rename(src, dst); err != nil {
		log.Println(err)
		return result{Info: "Error: failed file delete", Err: err}
	}
	return result{Info: "OK", Err: nil}
}

func infoStorageCmd(s storage) result {
	//qemu-img info [-f format] [filename]

	return RunStorageCmd([]string{"qemu-img", "info", s.path})
}
