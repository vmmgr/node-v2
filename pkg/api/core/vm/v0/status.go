package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"github.com/vmmgr/node/pkg/api/meta/json"
	"log"
	"net/http"
)

func (h *VMHandler) Startup(c *gin.Context) {
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

	if stat != libvirt.DOMAIN_RUNNING {
		if err := dom.Create(); err != nil {
			json.ResponseError(c, http.StatusInternalServerError, err)
		}
	}

	err = dom.Free()
	if err != nil {
		log.Println(err)
	}

	json.ResponseOK(c, nil)
}

func (h *VMHandler) Shutdown(c *gin.Context) {
	id := c.Param("id")

	var input vm.VirtualMachineStop

	err := c.BindJSON(&input)
	if err != nil {
		json.ResponseError(c, http.StatusBadRequest, err)
		return
	}
	log.Println(err)

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

	if stat != libvirt.DOMAIN_SHUTOFF {
		// Forceがtrueである場合、強制終了
		if input.Force {
			if err := dom.Destroy(); err != nil {
				json.ResponseError(c, http.StatusInternalServerError, err)
			}
		} else {
			if err := dom.Shutdown(); err != nil {
				json.ResponseError(c, http.StatusInternalServerError, err)
			}
		}
	}

	err = dom.Free()
	if err != nil {
		log.Println(err)
	}

	json.ResponseOK(c, nil)
}

func (h *VMHandler) Reset(c *gin.Context) {
	id := c.Param("id")

	dom, err := h.Conn.LookupDomainByUUIDString(id)
	if err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if err := dom.Reset(0); err != nil {
		json.ResponseError(c, http.StatusInternalServerError, err)
	}

	err = dom.Free()
	if err != nil {
		log.Println(err)
	}

	json.ResponseOK(c, nil)
}
