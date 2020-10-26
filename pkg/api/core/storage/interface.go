package storage

//Storage Type
//0: CDROM
//10: BootDevice(VirtIO)
//11: BootDevice(SATA)
//12: BootDevice(IDE)

type Storage struct {
	Type     uint   `json:"type"`
	FileType uint   `json:"filetype"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"readonly"`
	Boot     uint   `json:"boot"`
}

type GenerateStorageXml struct {
	Storage       Storage
	Number        uint
	AddressNumber uint
}

func GetDriverName(driver uint) string {
	// デフォルトはqcow2
	if driver == 1 {
		return "img"
	}
	return "qcow2"
}
