package snowflake

import (
	"errors"
	"net"
)

// PrivateIPToMachineID convert private ip to machine id.
// From https://github.com/sony/sonyflake/blob/master/sonyflake.go
func PrivateIPToMachineID() uint16 {
	ip, err := lower16BitPrivateIP()
	if err != nil {
		return 0
	}

	return ip
}

//--------------------------------------------------------------------
// private function defined.
//--------------------------------------------------------------------

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}

	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}
	return ipv4touint(ip), nil
}

func ipv4touint(ip net.IP) uint16 {
	l := len(ip)
	var r uint = uint(ip[l -1] & 0xFF)
	r |= uint(ip[l -2] ) << 8 & 0xFF00
	r |= uint(ip[l -3] ) << 16 & 0xFF0000
	r |= uint(ip[l -4] )<< 24 & 0xFF000000
	return uint16(r)
}
