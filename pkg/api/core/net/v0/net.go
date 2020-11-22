package v0

import (
	"github.com/vishvananda/netlink"
	"github.com/vmmgr/node/pkg/api/core/net"
)

type Net struct {
	name string
	link *netlink.Link
}

// ip link list
func ListSlaves() ([]net.Net, error) {
	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	var devices []net.Net
	for _, l := range links {
		devices = append(devices, net.Net{
			Name:  l.Attrs().Name,
			MAC:   l.Attrs().HardwareAddr.String(),
			MTU:   l.Attrs().MTU,
			Index: l.Attrs().Index,
		})
	}

	return devices, nil
}
