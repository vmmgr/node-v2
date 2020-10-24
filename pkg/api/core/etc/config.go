package etc

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Disk  DiskData  `json:"disk"`
	Image ImageData `json:image`
}

type DiskData struct {
	Path []DiskPath `json:"path"`
}

type DiskPath struct {
	Type   int    `json:"type"`
	Path   string `json:"path"`
	Status bool   `json:"status"`
}

type ImageData struct {
	Path string `json:"path"`
}

func GetDiskPath(id int) string {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Println(err)
	}
	var config Config
	json.Unmarshal(file, &config)
	for _, v := range config.Disk.Path {
		if v.Type == id && v.Status {
			return v.Path
		}
	}
	return ""
}

func GetImagePath() string {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println(err)
	}
	var config Config
	json.Unmarshal(file, &config)
	return config.Image.Path
}
