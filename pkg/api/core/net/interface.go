package net

type Net struct {
	Name  string `json:"name"`
	MAC   string `json:"mac"`
	MTU   int    `json:"mtu"`
	Index int    `json:"index"`
}
