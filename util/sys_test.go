package util

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
)

func TestIp(t *testing.T) {
	result := Ip2Long("98.138.253.109")
	fmt.Println(result)
	// or if you prefer the super fast way
	faster := binary.BigEndian.Uint32(net.ParseIP("98.138.253.109")[12:16])
	fmt.Println(faster)
	faster64 := int64(faster)
	fmt.Println(BackToIP4(faster64))
	ip1 := Ip2Long("38.0.0.0")
	ip2 := Ip2Long("38.0.0.255")
	//ip1 := Ip2Long("192.168.0.0")
	//ip2 := Ip2Long("192.168.0.255")
	x := ip2 - ip1
	fmt.Println(ip1, ip2, x)
	for i := ip1; i <= ip2; i++ {
		i := int64(i)
		fmt.Println(BackToIP4(i))
	}
}
func TestExternalIP(t *testing.T) {
	fmt.Println(ExternalIP())
	fmt.Println(GetMac())
}
