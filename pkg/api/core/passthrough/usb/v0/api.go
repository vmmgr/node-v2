package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/libvirt/libvirt-go"
	"github.com/vmmgr/node/pkg/api/meta/json"
	"log"
	"net/http"
)

type USBHandler struct {
	Conn *libvirt.Connect
}

func NewUSBHandler(handler USBHandler) *USBHandler {
	return &USBHandler{Conn: handler.Conn}
}

func (h *USBHandler) GetAPI(c *gin.Context) {
	dev, err := h.Get()
	if err != nil {
		log.Println(err)
		json.ResponseError(c, http.StatusInternalServerError, err)
	}
	json.ResponseOK(c, gin.H{
		"status": 1,
		"usb":    dev,
	})
}
