package etc

import (
	"fmt"
	"io"
	"os"
)

func FileCopy(srcName, dstName string) bool {
	fmt.Println("---CopyDisk")
	fmt.Println("src: " + srcName)
	fmt.Println("dst: " + dstName)
	src, err := os.Open(srcName)
	if err != nil {
		fmt.Println("OpenError")
		return false
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		fmt.Println("FIleCreateError")
		return false
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Println("CopyError")
		return false
	}
	return true
}
