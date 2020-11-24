package v0

import (
	"encoding/xml"
	"fmt"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"log"
	"sort"
)

var portCount = 500
var vncPortStart = 5910
var webSocketPortStart = 5310

type generateVNC struct {
	VNCPort       int
	WebSocketPort int
}

func (h *VMHandler) generateVNC() (generateVNC, error) {
	doms, err := h.Conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		log.Println(err)
	}

	//5900-6400
	var vncPort []int
	//5300-5800
	var webSocketPort []int

	for _, dom := range doms {
		t := libVirtXml.Domain{}
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &t)

		if len(t.Devices.Graphics) != 0 {
			for _, tmp := range t.Devices.Graphics {
				if tmp.VNC != nil {
					//5900-6400
					vncPort = append(vncPort, tmp.VNC.Port)
					//5300-5800
					webSocketPort = append(webSocketPort, tmp.VNC.WebSocket)
				}
			}
		}
	}

	//昇順に並び替える
	sort.Ints(vncPort)
	sort.Ints(webSocketPort)

	//ポート番号
	vncPortCount := vncPortStart
	webSocketPortCount := webSocketPortStart

	for _, port := range vncPort {
		//Port番号が上限に達する場合、エラーを返す
		if vncPortStart+portCount <= vncPortCount {
			return generateVNC{}, fmt.Errorf("Error: max port ")
		}

		if vncPortCount < port {
			break
		}
		vncPortCount++
	}

	for _, port := range webSocketPort {
		//Port番号が上限に達する場合、エラーを返す
		if webSocketPortCount+portCount <= webSocketPortCount {
			return generateVNC{}, fmt.Errorf("Error: max port ")
		}

		if webSocketPortCount < port {
			break
		}
		webSocketPortCount++
	}
	return generateVNC{VNCPort: vncPortCount, WebSocketPort: webSocketPortCount}, nil
}
