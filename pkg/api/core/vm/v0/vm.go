package v0

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	storageInterface "github.com/vmmgr/node/pkg/api/core/storage"
	storage "github.com/vmmgr/node/pkg/api/core/storage/v0"
	"github.com/vmmgr/node/pkg/api/core/tool/node"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"github.com/vmmgr/node/pkg/api/meta/json"
	"log"
	"net/http"
	"time"
)

type VMHandler struct {
	Conn *libvirt.Connect
	VM   vm.VirtualMachine
}

func NewVMHandler(input VMHandler) *VMHandler {
	return &VMHandler{Conn: input.Conn, VM: input.VM}
}

var process bool = false

func (h *VMHandler) Add(c *gin.Context) {

	if process {
		log.Println("Other VM Creating process... ")
		json.ResponseError(c, http.StatusInternalServerError, fmt.Errorf("Other VM Creating process... "))
		return
	}
	process = true

	//token1 := c.Request.Header.Get("TOKEN_1")
	//token2 := c.Request.Header.Get("TOKEN_2")

	var input vm.VirtualMachine

	err := c.BindJSON(&input)
	if err != nil {
		json.ResponseError(c, http.StatusBadRequest, err)
		process = false
		return
	}
	log.Println(err)

	// Templateで展開する場合
	if input.Template.Apply {
		log.Println(input.Template)
		storageh := storage.NewStorageHandler(storage.StorageHandler{Conn: h.Conn})

		go func() {
			err = storageh.Add(input.Template.Storage)
			if err != nil {
				log.Println(err)
				process = false
				node.SendServer(input.Info, 1, 1, "", err)
				return
			}
			timer := time.NewTimer(20 * time.Minute)
			defer timer.Stop()

			var err error

			//Todo 取りこぼす可能性があるので、要調査
		L:
			for {
				select {
				//20分以上かかる場合はタイムアウトさせる
				case <-timer.C:
					err = fmt.Errorf("Error: timeout ")

					break L
					//UUIDとGroupIDがMatchし、Progressが100の場合、storage転送処理が終了
				case msg := <-storageInterface.Path:
					//path変数にnode側のストレージをフルパスで代入する
					node.SendServer(input.Info, 2, 1, "Process: Creating VM", nil)
					input.Storage[0].Path = msg
					err = h.createVM(input)
					break L
				}
			}
			if err != nil {
				node.SendServer(input.Info, 2, 100, "Error: Create VM", err)
				return
			}
			process = false
			node.SendServer(input.Info, 2, 100, "Done: Create VM", nil)
			return
		}()

		json.ResponseError(c, http.StatusOK, nil)

	} else {
		//Template展開せず、手動でVMを作成する場合
		if err := h.createVM(input); err != nil {
			json.ResponseError(c, http.StatusInternalServerError, err)
		} else {
			json.ResponseError(c, http.StatusOK, nil)
		}
		process = false
		return
	}
}

func (h *VMHandler) Delete(c *gin.Context) {

	//token1 := c.Request.Header.Get("TOKEN_1")
	//token2 := c.Request.Header.Get("TOKEN_2")

	id := c.Param("id")

	dom, err := h.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}
	if stat == libvirt.DOMAIN_SHUTOFF {
		log.Println("power off")
	} else {
		if err := dom.Destroy(); err != nil {
			json.ResponseError(c, http.StatusInternalServerError, err)
			return
		}
	}

	if err := dom.Undefine(); err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	err = dom.Free()
	log.Println(err)

	json.ResponseOK(c, nil)
}

func (h *VMHandler) Update(c *gin.Context) {

}

func (h *VMHandler) Get(c *gin.Context) {
	id := c.Param("id")

	dom, err := h.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		log.Println(err)
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	// 初期定義
	t := libVirtXml.Domain{}

	// XMLをStructに代入
	tmpXml, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
	xml.Unmarshal([]byte(tmpXml), &t)

	json.ResponseOK(c, gin.H{
		"vm": vm.Detail{
			VM:   t,
			Stat: uint(stat),
		},
		"xml": tmpXml,
	})
}

func (h *VMHandler) GetStatus(c *gin.Context) {
	id := c.Param("id")

	dom, err := h.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	stat, _, err := dom.GetState()
	if err != nil {
		log.Println(err)
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	var status string

	if stat == libvirt.DOMAIN_SHUTOFF {
		log.Println("power off")
		status = "power off"
	} else if stat == libvirt.DOMAIN_RUNNING {
		status = "power on"
	} else if stat == libvirt.DOMAIN_SHUTDOWN {
		status = "shutdown"
	} else if stat == libvirt.DOMAIN_PAUSED {
		status = "paused"
	}

	err = dom.Free()
	log.Println(err)

	json.ResponseOK(c, gin.H{
		"status":     int(stat),
		"status_str": status,
	})
}

func (h *VMHandler) GetAll(c *gin.Context) {
	doms, err := h.Conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	log.Println(doms)
	if err != nil {
		log.Println(err)
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	var vms []vm.Detail

	for _, dom := range doms {
		t := libVirtXml.Domain{}
		stat, _, _ := dom.GetState()
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &t)

		//log.Println(len(t.Devices.Graphics))
		//log.Println(t.Devices.Graphics[0].VNC.Port)
		//log.Println(t.Devices.Graphics)
		vms = append(vms, vm.Detail{
			VM:   t,
			Stat: uint(stat),
		})
	}

	json.ResponseOK(c, gin.H{
		"vm": vms,
	})
}

func ChangeStatus(status bool) {
	process = status
}
