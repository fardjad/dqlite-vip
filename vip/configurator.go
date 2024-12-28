package vip

import (
	"time"

	"github.com/j-keck/arping"
	"github.com/vishvananda/netlink"
)

// Implements [Configurator]
type NetConfigurator struct {
	gratuitousArpCount int
}

// Implements [ConfiguratorFactoryFunc]
func NewNetConfigurator(gratuitousArpCount int) Configurator {
	return &NetConfigurator{gratuitousArpCount: gratuitousArpCount}
}

func (n *NetConfigurator) AddVIP(iface string, address string) error {
	parsedAddr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return err
	}

	err = netlink.AddrAdd(link, parsedAddr)
	if err != nil {
		return err
	}

	for i := 0; i < n.gratuitousArpCount; i++ {
		arping.GratuitousArpOverIfaceByName(parsedAddr.IP, iface)
		time.Sleep(1 * time.Second)
	}

	return err
}

func (n *NetConfigurator) RemoveVIP(iface string, address string) error {
	parsedAddr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return err
	}

	linkAddrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	for _, addr := range linkAddrs {
		if parsedAddr.Equal(addr) {
			return netlink.AddrDel(link, &addr)
		}
	}

	return nil
}
