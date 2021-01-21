package v0

import (
	"fmt"
	"github.com/vmmgr/node/pkg/api/core/vm"
	"log"
)

func (h *VMHandler) createVM(input vm.VirtualMachine) error {

	// VNC Portが0の場合、自動生成を行う
	if input.VNCPort == 0 {
		vnc, err := h.generateVNC()
		if err != nil {
			return err
		}
		input.VNCPort = uint(vnc.VNCPort)
		input.WebSocketPort = uint(vnc.WebSocketPort)
	}

	//メソッドにVM情報を代入
	h.VM = input

	domCfg, err := h.xmlGenerate()
	if err != nil {
		log.Println(err)
		return err
	}

	xml, err := domCfg.Marshal()
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(xml)

	dom, err := h.Conn.DomainDefineXML(xml)
	if err != nil {
		return err
	}

	err = dom.Create()
	if err != nil {
		// node側でエラーを表示
		log.Println(err)
		return err
	}
	return nil
}
