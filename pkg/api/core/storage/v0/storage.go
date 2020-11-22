package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/node/pkg/api/core/storage"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
	"github.com/vmmgr/node/pkg/api/meta/json"
	"log"
	"net/http"
)

type StorageHandler struct {
	conn *libvirt.Connect
}

func NewStorageHandler(connect *libvirt.Connect) *StorageHandler {
	return &StorageHandler{conn: connect}
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
			path = tmpConf.Path + "/" + input.Path
		}
	}

	log.Println(path)
	// Pathが見つからない場合
	if path == "" {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: Not found... "))
		return
	}

	if fileExistsCheck(path) {
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
			log.Println(input.FromImaCon.Path)
			log.Println(input.Path)
			err := sftpRemoteToLocal(storage.SFTPAuth{
				IP: input.FromImaCon.IP, User: config.Conf.ImaCon.User, Pass: config.Conf.ImaCon.Pass,
			}, input.FromImaCon.Path, path, input)
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
	if !fileExistsCheck(input.SrcFile) {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: file no exists... "))
		return
	}

	// Destinationファイルの確認
	if fileExistsCheck(input.DstFile) {
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
	if !fileExistsCheck(input.SrcFile) {
		json.ResponseError(c, http.StatusNotFound, fmt.Errorf("Error: file no exists... "))
		return
	}

	if data, err := infoImage(input.SrcFile); err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
	} else {
		json.ResponseOK(c, data)
	}
}
