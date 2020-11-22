package v0

import (
	"github.com/vishvananda/netlink"
)

type Bond struct {
	name string
	link *netlink.Bond
}

// ip link add bond0 type bond miimon 100 mode 802_3ad
func (b *Bond) createBond() error {
	l := netlink.NewLinkAttrs()
	l.Name = b.name
	b.link = &netlink.Bond{
		LinkAttrs:      netlink.LinkAttrs{Name: b.name, HardwareAddr: nil},
		Mode:           netlink.BOND_MODE_802_3AD,
		XmitHashPolicy: netlink.BOND_XMIT_HASH_POLICY_LAYER2_3,
		Miimon:         100,
	}
	if err := netlink.LinkAdd(b.link); err != nil {
		return err
	}
	if err := b.Up(); err != nil {
		return err
	}
	return nil
}

func (b *Bond) Up() error {
	if err := netlink.LinkSetUp(b.link); err != nil {
		return err
	}
	return nil
}
