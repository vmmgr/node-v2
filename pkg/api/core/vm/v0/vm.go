package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"github.com/vmmgr/node/pkg/api/meta/json"
	"log"
	"net/http"
)

type VMHandler struct {
	conn *libvirt.Connect
}

func NewMainHandler(connect *libvirt.Connect) *VMHandler {
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

		info, _ := domResult.GetInfo()
		uuid, _ := domResult.GetUUIDString()
		_, stat, _ := domResult.GetState()
		vms = append(vms, vm.VirtualMachine{
			Name:   dom,
			Memory: uint(info.Memory),
			UUID:   uuid,
			Stat:   uint(stat),
		})
		fmt.Println(info)
		fmt.Println(uuid)
	}
	json.ResponseOK(c, gin.H{
		"vm": vms,
	})
}
