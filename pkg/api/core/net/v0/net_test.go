package v0

import (
	"encoding/xml"
	"fmt"
	"github.com/libvirt/libvirt-go"
	libVirtXml "github.com/libvirt/libvirt-go-xml"
	"github.com/vmmgr/node/pkg/api/core/tool/config"
	"sort"
	"strconv"
	"strings"
	"testing"
)

var maxMAC = 65535
var startMAC = 10

func Test_NetList(t *testing.T) {
	var a []int

	for i := 0; i < 65535; i++ {
		a = append(a, i)
	}

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		t.Fatal("failed to connect to qemu")
	}
	defer conn.Close()

	var macs []int

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	t.Log(doms)
	if err != nil {
		t.Fatal(err)
	}

	// Todo:
	for _, dom := range doms {
		data := libVirtXml.Domain{}
		xmlString, _ := dom.GetXMLDesc(libvirt.DOMAIN_XML_SECURE)
		xml.Unmarshal([]byte(xmlString), &data)

		if len(data.Devices.Interfaces) != 0 {
			for _, tmp := range data.Devices.Interfaces {
				t.Log(tmp.MAC.Address)
				mac := strings.Split(tmp.MAC.Address, ":")
				if (mac[0] + mac[1]) == "5254" {
					v, _ := strconv.ParseInt(mac[4]+mac[5], 16, 0)
					macs = append(macs, int(v))
					t.Log(v)
				}
			}
		}
	}

	//昇順に並び替える
	sort.Ints(macs)

	//startMACを定義
	mac := startMAC

	for _, m := range macs {
		//Port番号が上限に達する場合、エラーを返す
		if maxMAC <= mac {
			t.Fatal("Error: max mac address")
		}
		if mac < m {
			break
		}
		mac++
	}

	mac = 255
	macIndex1 := mac / 256
	macIndex2 := mac % 256

	t.Log(fmt.Sprintf("%s:%.2x:%.2x", config.Conf.Node.MAC, macIndex1, macIndex2))

	t.Log(config.Conf.Node.MAC)
}
