package vip

import (
	"log"
	"time"

	"github.com/j-keck/arping"
	"github.com/vishvananda/netlink"
)

// Implements [Configurator]
type configurator struct {
	gratuitousArpCount int
}

// Implements [ConfiguratorFactoryFunc]
func NewConfigurator(gratuitousArpCount int) Configurator {
	return &configurator{gratuitousArpCount: gratuitousArpCount}
}

func (n *configurator) hasVIP(link netlink.Link, addrToCheck *netlink.Addr) (bool, error) {
	linkAddrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return false, err
	}

	for _, addr := range linkAddrs {
		if addr.Equal(*addrToCheck) {
			return true, nil
		}
	}

	return false, nil
}

func (n *configurator) EnsureVIP(iface string, address string) error {
	parsedAddr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return err
	}

	if hasVIP, err := n.hasVIP(link, parsedAddr); err != nil {
		return err
	} else if hasVIP {
		return nil
	}

	err = netlink.AddrAdd(link, parsedAddr)
	if err != nil {
		return err
	}

	log.Printf("Added VIP %s to interface %s", address, iface)

	for i := 0; i < n.gratuitousArpCount; i++ {
		arping.GratuitousArpOverIfaceByName(parsedAddr.IP, iface)
		time.Sleep(1 * time.Second)
	}

	return err
}

func (n *configurator) RemoveVIP(iface string, address string) error {
	parsedAddr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return err
	}

	if hasVIP, err := n.hasVIP(link, parsedAddr); err != nil {
		return err
	} else if !hasVIP {
		return nil
	}

	return netlink.AddrDel(link, parsedAddr)
}
