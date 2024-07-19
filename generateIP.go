package main 

import (
	"log"
	"net"
)

func ipMask() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return ""
	}
	
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		
		addrs, err := iface.Addrs()
		if err != nil {
			log.Println(err)
			continue
		}
		
		for _, addr := range addrs {
			ip, ipnet, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.Println(err)
				continue
			}
			
			if ip.To4() == nil {
				continue
			}
			
			network := ip.Mask(ipnet.Mask)
			networkCIDR := &net.IPNet{IP: network, Mask: ipnet.Mask}
			return networkCIDR.String()
		}
	}
	
	return ""
}
