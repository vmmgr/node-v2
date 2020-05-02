package etc

import (
	"fmt"
	"github.com/yoneyan/vm_mgr/node/db"
	"math/rand"
	"time"
)

const (
	sockpath = "/kvm/socket"
)

func SocketConnectionPath(socketfile string) string {
	socket := "unix-connect:" + sockpath + "/" + socketfile + ".sock"
	return socket
}

func SocketPath(socketfile string) string {
	return sockpath + "/" + socketfile + ".sock"
}

func SocketGenerate(socketfile string) string {
	//unix:/tmp/monitor.sock,server,nowait
	return "unix:" + sockpath + "/" + socketfile + ".sock,server,nowait"
}

func GeneratePath(path, name string) string {
	return path + "/" + name
}

func RandomGenerateValue() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(600)
}

func GenerateVNCPort() int {
	var vncarray []int
	data := db.VMDBGetAll()
	for _, a := range data {
		vncarray = append(vncarray, a.Vnc)
	}
	fmt.Println(vncarray)

	var vnc int
	for {
		vnc = RandomGenerateValue()
		vnc = vnc + 200
		if valuecontain(vncarray, vnc) == false {
			break
		}
	}
	fmt.Print("vnc: ")
	fmt.Println(vnc)
	return vnc
}

func valuecontain(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
