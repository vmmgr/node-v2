package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Controller Controller `json:"controller"`
	Node       Node       `json:"node"`
	ImaCon     ImaCon     `json:"imacon"`
	Storage    []Storage  `json:"storage"`
}

type Controller struct {
	List []List `json:"port"`
	Auth Auth   `json:"auth"`
}

type Auth struct {
	Token1 string `json:"token1"`
	Token2 string `json:"token2"`
	Token3 string `json:"token3"`
}

type List struct {
	URL string `json:"url"`
}

type Node struct {
	IP      string `json:"ip"`
	Port    uint   `json:"port"`
	Machine string `json:"machine"`
}

type ImaCon struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Storage struct {
	ID   uint   `json:"id"`
	Type uint   `json:"type"`
	Path string `json:"path"`
}

var Conf Config

func GetConfig(inputConfPath string) error {
	configPath := "./data.json"
	if inputConfPath != "" {
		configPath = "./" + inputConfPath
	}
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	var data Config
	json.Unmarshal(file, &data)
	Conf = data
	return nil
}
