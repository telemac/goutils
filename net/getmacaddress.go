package net

import (
	"errors"
	"net"
)

type networkInterface struct {
	Name      string
	IPAddress string
	MAC       string
}

func LocalAddresses() ([]networkInterface, error) {
	var networkInterfaces []networkInterface

	ifaces, err := net.Interfaces()
	if err != nil {
		return networkInterfaces, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			// log.Print(Errorf("localAddresses error : %v\n", err.Error()))
			continue
		}
		for _, a := range addrs {

			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					interface_name := i.Name
					// fmt.Println("interface_name =", i )
					interface_ipAddress := a.(*net.IPNet).IP.String()
					// Println("interface_ipAddress =", interface_ipAddress )
					networkInterfaces = append(networkInterfaces, networkInterface{Name: interface_name, IPAddress: interface_ipAddress, MAC: i.HardwareAddr.String()})
				}
			}
		}
	}
	//	log.Println("networkInterfaces",networkInterfaces)
	return networkInterfaces, nil
}

// GetMACAddress gets the mac address of the first non loopback ipv4 interface
func GetMACAddress() (string, error) {
	NetworkInterfaces, err := LocalAddresses()
	if err != nil {
		return "", err
	}
	if len(NetworkInterfaces) == 0 {
		return "", errors.New("GetMACAddress: no interface detected")
	}
	return NetworkInterfaces[0].MAC, nil
}
