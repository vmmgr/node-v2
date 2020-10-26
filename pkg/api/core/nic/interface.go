package nic

const (
	s = 0
)

//Type
//0: Bridge 1: NAT 2:macvtap
//Mode
//0: Bridge 1: vpea 2: private 3: passthrough

type NIC struct {
	Type   uint   `json:"type"`
	Driver uint   `json:"driver"`
	Mode   uint   `json:"mode"`
	MAC    string `json:"mac"`
	Device string `json:"device"`
}

type GenerateNICXml struct {
	NIC           NIC
	AddressNumber uint
}

func GetDriverName(driver uint) string {
	//デフォルトはvirtio

	if driver == 0 {
		return "e1000e"
	} else if driver == 1 {
		return "rtl8139"
	}
	return "virtio"
}

//Mode
func GetModeName(mode uint) string {
	//デフォルトはBrdige

	if mode == 1 {
		return "vepa"
	} else if mode == 2 {
		return "private"
	} else if mode == 3 {
		return "passthrough"
	}
	return "bridge"
}
