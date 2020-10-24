package etc

import (
	"strconv"
)

const (
	sockPath = "/kvm/socket"
)

func SocketConnectionPath(id int) string {
	socket := "unix-connect:" + sockPath + "/" + strconv.Itoa(id) + ".sock"
	return socket
}

func SocketPath(id int) string {
	return sockPath + "/" + strconv.Itoa(id) + ".sock"
}

func SocketGenerate(id int) string {
	//unix:/tmp/monitor.sock,server,nowait
	return "unix:" + sockPath + "/" + strconv.Itoa(id) + ".sock,server,nowait"
}
