package v0

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar"
	"github.com/vmmgr/controller/pkg/api/core/controller"
	"github.com/vmmgr/node/pkg/api/core/storage"
	"github.com/vmmgr/node/pkg/api/core/tool/client"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"time"
)

type File struct {
	uuid string
}

type Progress struct {
	total int64
	size  int64
}

func (p *Progress) Write(data []byte) (int, error) {
	n := len(data)
	p.size += int64(n)

	return n, nil
}

func (h *StorageHandler) sftpRemoteToLocal() error {
	//config := &ssh.ClientConfig{User: auth.User, HostKeyCallback: nil, Auth: []ssh.AuthMethod{ssh.Password(auth.Pass)}}
	config := &ssh.ClientConfig{User: h.Auth.User, HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{ssh.Password(h.Auth.Pass)}}
	config.SetDefaults()
	sshConn, err := ssh.Dial("tcp", h.Auth.IP+":22", config)
	if err != nil {
		log.Println(err)
		sendServer(h.Input, h.DstPath, 0, err)
		return err
	}
	defer sshConn.Close()

	// SFTP Client
	client, err := sftp.NewClient(sshConn)
	if err != nil {
		log.Println(err)
		sendServer(h.Input, h.DstPath, 0, err)
		return err
	}
	defer client.Close()

	// dstFileの作成
	dstFile, err := os.Create(h.DstPath)
	if err != nil {
		log.Println(err)
		sendServer(h.Input, h.DstPath, 0, err)
		return err
	}
	defer dstFile.Close()

	// srcFileをOpen
	srcFile, err := client.Open(h.SrcPath)
	if err != nil {
		log.Println(err)
		sendServer(h.Input, h.DstPath, 0, err)
		return err
	}

	file, err := srcFile.Stat()
	if err != nil {
		log.Println("Error: file gateway error")
		sendServer(h.Input, h.DstPath, 0, err)
		return err
	}

	log.Println(file.Size())

	p := Progress{total: file.Size()}

	count := 100
	count64 := int64(count)
	bar := progressbar.Default(count64)

	// Node側の表示
	go func() {
		for {
			if p.size != p.total {
				<-time.NewTimer(200 * time.Microsecond).C
				bar.Set(int(float64(p.size) / float64(p.total) * 100))
			} else {
				return
			}
		}
	}()

	// Node側の表示
	go func() {
		for {
			if p.size != p.total {
				<-time.NewTimer(1 * time.Second).C
				sendServer(h.Input, h.DstPath, int(float64(p.size)/float64(p.total)*100), nil)
			} else {
				return
			}
		}
	}()

	// コピーの処理
	bytes, err := io.Copy(dstFile, io.TeeReader(srcFile, &p))
	if err != nil {
		log.Println(err)
		sendServer(h.Input, h.DstPath, 0, err)
		return err
	}
	bar.Set(100)
	fmt.Printf("\n%dbytes copied\n", bytes)

	// sync
	err = dstFile.Sync()
	if err != nil {
		log.Println(err)
		sendServer(h.Input, h.DstPath, 0, err)
		return err
	}
	_, err = capacityExpansion(h.DstPath, h.Input.Capacity)
	if err != nil {
		log.Println("Error: disk capacity expansion")
		log.Println(err)
	}

	sendServer(h.Input, h.DstPath, 100, nil)

	return nil
}

func sftpLocalToRemote(auth storage.SFTPAuth, srcLocalPath, dstRemotePath string) {
	//config := &ssh.ClientConfig{User: auth.User, HostKeyCallback: nil, Auth: []ssh.AuthMethod{ssh.Password(auth.Pass)}}
	config := &ssh.ClientConfig{User: auth.User, HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{ssh.Password(auth.Pass)}}
	config.SetDefaults()
	sshConn, err := ssh.Dial("tcp", "example.com:22", config)
	if err != nil {
		panic(err)
	}
	defer sshConn.Close()

	// SFTP Client
	client, err := sftp.NewClient(sshConn)
	if err != nil {
		log.Println(err)
	}
	defer client.Close()

	// dstFileの作成
	dstFile, err := client.Create(dstRemotePath)
	if err != nil {
		log.Println(err)
	}
	defer dstFile.Close()

	srcFile, err := os.Open(srcLocalPath)
	if err != nil {
		log.Println(err)
	}

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%d bytes copied\n", bytes)
}

func fileCopy(srcFile, dstFile, controller string) error {
	log.Println("---Copy disk image")
	log.Println("src: " + srcFile)
	log.Println("dst: " + dstFile)
	src, err := os.Open(srcFile)
	if err != nil {
		log.Println("Error: open error")
		return fmt.Errorf("open error")
	}
	defer src.Close()
	file, err := src.Stat()
	if err != nil {
		log.Println("Error: file gateway error")
		return err
	}

	dst, err := os.Create(dstFile)
	if err != nil {
		log.Println("Error: file create")
		return err
	}
	defer dst.Close()

	p := Progress{total: file.Size()}

	count := 100
	count64 := int64(count)
	bar := progressbar.Default(count64)

	go func() {
		for {
			if p.size != p.total {
				<-time.NewTimer(200 * time.Microsecond).C
				//log.Println(tmp.fileSize)
				bar.Set(int(float64(p.size) / float64(p.total) * 100))
				//sendServer()
			} else {
				log.Println("end")
				return
			}
		}
	}()

	_, err = io.Copy(dst, io.TeeReader(src, &p))
	if err != nil {
		log.Println("Error: file copy error")
		return err
	}

	return nil
}

func sendServer(input storage.Storage, filePath string, progress int, error error) {
	for _, srv := range config.Conf.Controller.List {
		sendBody, _ := json.Marshal(controller.Node{
			GroupID:  input.GroupID,
			UUID:     input.UUID,
			FilePath: filePath,
			Progress: uint(progress),
			Error:    error,
			Comment:  "storage creating...",
		})
		client.Post(srv.URL, sendBody)
	}
}
