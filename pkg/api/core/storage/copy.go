package storage

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type FileTransfer struct {
	io.Reader
	total    int64
	fileSize int64
}

func (ft *FileTransfer) Read(p []byte) (int, error) {
	n, err := ft.Reader.Read(p)
	ft.total += int64(n)
	return n, err
}

func (t *Tmp) fileCopy(srcFile, dstFile string) result {
	log.Println("---Copy disk image")
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
		log.Println("Error: file gateway error")
		return result{Info: "Error: file create", Err: err}
	}

	dst, err := os.Create(dstFile)
	if err != nil {
		log.Println("Error: file create")
		return result{Info: "Error: file create", Err: err}
	}
	defer dst.Close()

	tmp := FileTransfer{Reader: src, fileSize: file.Size()}
	copy := true
	go func() {
		for tmp.fileSize != tmp.total && copy == true {
			timer := time.NewTimer(time.Second * 1)
			<-timer.C
			t.Info = "Progress: " + strconv.Itoa(int(tmp.total)) + " FileSize: " + strconv.Itoa(int(tmp.fileSize))
		}
	}()

	_, err = io.Copy(dst, &tmp)
	if err != nil {
		log.Println("Error: file copy error")
		copy = false
		return result{Info: "Error: file copy error", Err: err}
	}
	return result{Info: "OK", Err: nil}
}
