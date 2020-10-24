package nic

import (
	"log"
	"math/rand"
	"strconv"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func generateMacAddress() string {
	mac := "52:54"
	var value string
	for i := 0; i < 4; i++ {
		value = strconv.FormatInt(int64(random(0, 255)), 16)
		if len(value) == 1 {
			mac += ":0" + value
		} else {
			mac += ":" + value
		}
	}
	log.Println("Generate: MAC Address: " + mac)
	return mac
}
