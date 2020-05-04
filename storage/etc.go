package storage

type storageDriver struct {
	driver    int
	extension string
}

func getDriver(driver int) storageDriver {
	extension := "qcow2"

	if driver == 2 {
		extension = "img"
	}
	return storageDriver{driver: driver, extension: extension}
}
