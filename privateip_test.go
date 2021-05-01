package snowflake

import (
	"net"
	"testing"
)

func TestPrivateIPToMachineID(t *testing.T) {
	mid := PrivateIPToMachineID()
	if mid <= 0 {
		t.Error("MachineID should be > 0")
	}
}

func TestPrivateIPToUintD(t *testing.T) {
	t.Log(privateIPv4())
	t.Log(lower16BitPrivateIP())
	t.Log(ipv4touint(net.IPv4(byte(10),byte(67),byte(4),byte(11))))
	t.Log(ipv4touint(net.IPv4(byte(127),byte(0),byte(0),byte(1))))
	t.Log(ipv4touint(net.IPv4(byte(10),byte(67),byte(5),byte(10))))
	t.Log(ipv4touint(net.IPv4(byte(10),byte(67),byte(4),byte(16))))
	t.Log(ipv4touint(net.IPv4(byte(10),byte(67),byte(4),byte(255))))
	t.Log(ipv4touint(net.IPv4(byte(255),byte(255),byte(255),byte(255))))

}
