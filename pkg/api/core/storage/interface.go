package storage

import "github.com/vmmgr/node/pkg/api/core/gateway"

var Path = make(chan string)

type Storage struct {
	Info       gateway.Info `json:"info"`           //Info
	Mode       uint         `json:"modestorageTmp"` //0:Manual 1:From ImaCon
	FromImaCon ImaCon       `json:"from_imacon"`    //Imageをpullする際に使用するURL
	VMName     string       `json:"vm_name"`        //VMNameがある場合は、Pathに追加する
	Type       uint         `json:"type"`           //1: CDROM 2:Floppy (no support) 10:BootDev(VirtIO) 11: BootDev(SATA) 12: BootDev(IDE)
	FileType   uint         `json:"filetype"`       //0:qcow2 1:img
	PathType   uint         `json:"path_type"`      //node側のストレージの種類 0~9:SSD 10~19:HDD 20~29:NVMe 100~109:SSD(iSCSI) 110~119:SSD(iSCSI) 120~129:NVme(iSCSI)
	Path       string       `json:"path"`           //node側のパス
	Capacity   uint         `json:"capacity"`       //容量
	ReadOnly   bool         `json:"readonly"`       //Readonlyであるか
	Boot       uint         `json:"boot"`
}

type VMStorage struct {
	Type     uint   `json:"type"`      //0:BootDev(VirtIO) 1: CDROM 2:Floppy (no support) 11: BootDev(SATA) 12: BootDev(IDE)
	FileType uint   `json:"filetype"`  //0:qcow2 1:img
	PathType uint   `json:"path_type"` //node側のストレージの種類 0~9:SSD 10~19:HDD 20~29:NVMe 100~109:SSD(iSCSI) 110~119:SSD(iSCSI) 120~129:NVme(iSCSI)
	Path     string `json:"path"`      //node側のパス
	ReadOnly bool   `json:"readonly"`  //Readonlyであるか
	Boot     uint   `json:"boot"`
}

type ImaCon struct {
	IP   string `json:"ip"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Path string `json:"path"`
}

type Convert struct {
	SrcFile string `json:"src_file"`
	SrcType string `json:"src_type"`
	DstFile string `json:"dst_file"`
	DstType string `json:"dst_type"`
}

type GenerateStorageXml struct {
	Storage       VMStorage
	Number        uint
	AddressNumber uint
}

type FileTransfer struct {
	URL         string
	CurrentSize int64
	AllSize     int64
}

type SFTPAuth struct {
	IP   string
	User string
	Pass string
}

func GetExtensionName(extension uint) string {
	// デフォルトはqcow2
	if extension == 1 {
		return "img"
	}
	return "qcow2"
}
