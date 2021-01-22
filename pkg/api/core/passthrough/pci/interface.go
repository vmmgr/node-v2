package pci

type PCI struct {
	ID         string `json:"id"`
	ClassName  string `json:"class_name"`
	VendorID   string `json:"Vendor_id"`
	VendorName string `json:"Vendor_name"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	Comment    string `json:"comment"`
}
