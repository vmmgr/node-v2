package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/imacon/pkg/api/meta/json"
	"log"
	"net/http"
)

type PCIHandler struct {
	Conn *libvirt.Connect
}

func NewPCIHandler(handler PCIHandler) *PCIHandler {
	return &PCIHandler{Conn: handler.Conn}
}

func (h *PCIHandler) GetAPI(c *gin.Context) {
	dev, err := h.Get()
	if err != nil {
		log.Println(err)
		json.ResponseError(c, http.StatusInternalServerError, err)
	}
	json.ResponseOK(c, gin.H{
		"status": 1,
		"pci":    dev,
	})
}
