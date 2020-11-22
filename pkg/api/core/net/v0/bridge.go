package v0

import "github.com/vishvananda/netlink"

type Bridge struct {
	name string
	link *netlink.Bridge
}

// ip link add $name type bridge
func (b *Bridge) createBridge() error {
	l := netlink.NewLinkAttrs()
	l.Name = b.name
	b.link = &netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: l.Name, HardwareAddr: nil}}
	if err := netlink.LinkAdd(b.link); err != nil {
		return err
	}
	if err := b.Up(); err != nil {
		return err
	}

	return nil
}

func (b *Bridge) Up() error {
	if err := netlink.LinkSetUp(b.link); err != nil {
		return err
	}
	return nil
}
