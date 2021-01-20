package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/node/pkg/api/core/storage"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"github.com/vmmgr/node/pkg/api/meta/json"
	"log"
	"net/http"
	"os"
)

type StorageHandler struct {
	Conn    *libvirt.Connect
	Input   storage.Storage
	VM      vm.VirtualMachine
	Address *vm.Address
	Auth    *storage.SFTPAuth
	SrcPath string
	DstPath string
}

func NewStorageHandler(handler StorageHandler) *StorageHandler {
	return &StorageHandler{Conn: handler.Conn, Input: handler.Input, VM: handler.VM, Address: handler.Address,
		Auth: handler.Auth, SrcPath: handler.SrcPath, DstPath: handler.DstPath}
}

func (h *StorageHandler) Add(c *gin.Context) {
	var input storage.Storage

	err := c.BindJSON(&input)
	if err != nil {
		json.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	path := ""

	for _, tmpConf := range config.Conf.Storage {
		if tmpConf.Type == input.PathType {
			if input.VMName == "" {
				path = tmpConf.Path + "/" + input.Path
			} else {
				if err := os.Mkdir(tmpConf.Path+"/"+input.VMName, 0775); err != nil {
					log.Println(err)
					json.ResponseError(c, http.StatusInternalServerError, err)
					return
				}
				path = tmpConf.Path + "/" + input.VMName + "/" + input.Path
			}
		}
	}

	log.Println(path)
	// Pathが見つからない場合
	if path == "" {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: Not found... "))
		return
	}

	if FileExistsCheck(path) {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: file already exists... "))
		return
	}

	var out string

	// イメージの作成
	if input.Mode == 0 {
		out, err = generateImage(storage.GetExtensionName(input.Type), input.Path, input.Capacity)
		if err != nil {
			json.ResponseError(c, http.StatusNotFound, err)
			return
		} else {
			json.ResponseOK(c, out)
		}
	} else if input.Mode == 1 {
		// ImaConからイメージ取得(時間がかかるので、go funcにて処理)
		go func() {
			log.Println("From: " + input.FromImaCon.Path)
			log.Println("To: " + path)

			//メソッドに各種情報の追加
			h.Auth = &storage.SFTPAuth{
				IP: input.FromImaCon.IP, User: config.Conf.ImaCon.User, Pass: config.Conf.ImaCon.Pass,
			}
			h.SrcPath = input.FromImaCon.Path
			h.DstPath = path
			h.Input = input

			err := h.sftpRemoteToLocal()
			log.Println(err)
		}()

		json.ResponseOK(c, out)
	}
}

func (h *StorageHandler) ConvertImage(c *gin.Context) {
	var input storage.Convert

	err := c.BindJSON(&input)
	if err != nil {
		json.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	// sourceファイルの確認
	if !FileExistsCheck(input.SrcFile) {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: file no exists... "))
		return
	}

	// Destinationファイルの確認
	if FileExistsCheck(input.DstFile) {
		json.ResponseError(c, http.StatusInternalServerError, fmt.Errorf("Error: file already exists... "))
		return
	}

	if err := convertImage(input); err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
	} else {
		json.ResponseOK(c, nil)
	}
}

func (h *StorageHandler) InfoImage(c *gin.Context) {
	var input storage.Convert

	err := c.BindJSON(&input)
	if err != nil {
		json.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	// sourceファイルの確認
	if !FileExistsCheck(input.SrcFile) {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: file no exists... "))
		return
	}

	if data, err := infoImage(input.SrcFile); err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
	} else {
		json.ResponseOK(c, data)
	}
}
