package storage

import (
	"io"
	"log"
	"os"
)

type fileTransfer struct {
	io.Reader
	total    int64
	fileSize int64
}

func (ft *fileTransfer) Read(p []byte) (int, error) {
	n, err := ft.Reader.Read(p)
	ft.total += int64(n)

	if err == nil {
		log.Println("Progress: ", ft.total, " FileSize: ", ft.fileSize)
	}
	return n, err
}

func fileCopy(srcFile, dstFile string) result {
	log.Println("---Copy Image")
	log.Println("src: " + srcFile)
	log.Println("dst: " + dstFile)
	src, err := os.Open(srcFile)
	if err != nil {
		log.Println("Error: open error")
		return result{Info: "Error: open error", Err: err}
	}
	defer src.Close()
	file, err := src.Stat()
	if err != nil {
		log.Println("Error: file data error")
		return result{Info: "Error: file create", Err: err}
	}

	dst, err := os.Create(dstFile)
	if err != nil {
		log.Println("Error: file create")
		return result{Info: "Error: file create", Err: err}
	}
	defer dst.Close()

	tmp := fileTransfer{Reader: src, fileSize: file.Size()}

	_, err = io.Copy(dst, &tmp)
	if err != nil {
		log.Println("Error: file copy error")
		return result{Info: "Error: file copy error", Err: err}
	}
	return result{Info: "OK", Err: nil}
}
