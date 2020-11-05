package v0

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"github.com/vmmgr/node/pkg/api/meta/json"
	"log"
	"net/http"
)

type VMHandler struct {
	conn *libvirt.Connect
}

func NewVMHandler(connect *libvirt.Connect) *VMHandler {
	return &VMHandler{conn: connect}
}

func (h *VMHandler) Add(c *gin.Context) {

	//token1 := c.Request.Header.Get("TOKEN_1")
	//token2 := c.Request.Header.Get("TOKEN_2")

	var input vm.VirtualMachine

	err := c.BindJSON(&input)
	if err != nil {
		json.ResponseError(c, http.StatusBadRequest, err)
		return
	}
	log.Println(err)

	domCfg, err := xmlGenerate(input)
	if err != nil {
		log.Println(err)
		json.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	xml, err := domCfg.Marshal()
	if err != nil {
		log.Println(err)
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	fmt.Println(xml)

	_, err = h.conn.DomainDefineXML(xml)
	if err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}
	//if err = dom.Create(); err != nil {
	//	json.ResponseError(c, http.StatusInternalServerError, err)
	//} else {
	json.ResponseOK(c, nil)
	//meta.ResponseJSON(c, http.StatusOK, nil, nil)
	//}
}

func (h *VMHandler) Delete(c *gin.Context) {

	//token1 := c.Request.Header.Get("TOKEN_1")
	//token2 := c.Request.Header.Get("TOKEN_2")

	id := c.Param("id")

	dom, err := h.conn.LookupDomainByUUIDString(id)
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

	dom, err := h.conn.LookupDomainByUUIDString(id)
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
		"vm": vm.VirtualMachine{
			Name:   t.Name,
			VCPU:   t.VCPU.Value,
			Memory: t.Memory.Value,
			UUID:   t.UUID,
			OS: vm.OS{
				Boot: t.OS.BootDevices[0].Dev,
				Type: t.OS.Type.Type,
			},
			Stat: uint(stat),
		},
		"xml": tmpXml,
	})
}

func (h *VMHandler) GetStatus(c *gin.Context) {
	id := c.Param("id")

	dom, err := h.conn.LookupDomainByUUIDString(id)
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
	doms, err := h.conn.ListDefinedDomains()
	if err != nil {
		log.Println(err)
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	var vms []vm.VirtualMachine

	for _, dom := range doms {
		log.Println(dom)
		domResult, err := h.conn.LookupDomainByName(dom)
		if err != nil {
			log.Println(err)
			json.ResponseError(c, http.StatusInternalServerError, err)
			return
		}

		t := libVirtXml.Domain{}

		_, stat, _ := domResult.GetState()
		tmp, _ := domResult.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(tmp), &t)

		//log.Println(t.OS.Type.Arch)
		vms = append(vms, vm.VirtualMachine{
			Name:   dom,
			VCPU:   t.VCPU.Value,
			Memory: t.Memory.Value,
			UUID:   t.UUID,
			OS: vm.OS{
				Boot: t.OS.BootDevices[0].Dev,
				Type: t.OS.Type.Type,
			},
			Stat: uint(stat),
		})
	}
	json.ResponseOK(c, gin.H{
		"vm": vms,
	})
}
