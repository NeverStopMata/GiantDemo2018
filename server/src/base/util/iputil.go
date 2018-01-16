package util

import (
	"net"
	"strings"
)

const (
	big     = 0xFFFFFF
	IPv4len = 4
)

var (
	inners []net.IP
	outers []net.IP
)

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	infiler1 := net.IPNet{IP: net.IPv4(10, 0, 0, 0), Mask: net.IPv4Mask(255, 0, 0, 0)}
	infiler2 := net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.IPv4Mask(255, 240, 0, 0)}
	infiler3 := net.IPNet{IP: net.IPv4(192, 168, 0, 0), Mask: net.IPv4Mask(255, 255, 0, 0)}
	infiler4 := net.IPNet{IP: net.IPv4(100, 64, 0, 0), Mask: net.IPv4Mask(255, 192, 0, 0)}
	selfIP := net.IPv4(127, 0, 0, 1)

	for _, addr := range addrs {
		var ip4 net.IP
		switch ipaddr := addr.(type) {
		case *net.IPNet:
			ip4 = ipaddr.IP.To4()
		case *net.IPAddr:
			ip4 = ipaddr.IP.To4()
		}
		if ip4 != nil && !ip4.Equal(selfIP) {
			if infiler1.Contains(ip4) || infiler2.Contains(ip4) || infiler3.Contains(ip4) || infiler4.Contains(ip4) {
				inners = append(inners, ip4)
			} else {
				outers = append(outers, ip4)
			}
		}
	}
}

func GetInnerIP() []net.IP {
	return inners
}

func GetTopInnerIP() net.IP {
	if len(inners) == 0 {
		return nil
	}
	return inners[0]
}

func GetOuterIP() []net.IP {
	return outers
}

func GetOuterIPStr() (str string) {
	for _, ip := range outers {
		if len(str) > 0 {
			str += "/"
		}
		str += ip.String()
	}
	return
}

func GetTopOuterIP() net.IP {
	if len(outers) == 0 {
		return nil
	}
	return outers[0]
}

func IPStrToUInt(str string) (n uint32) {
	strs := strings.Split(str, ".")
	if len(strs) != 4 {
		return 0
	}
	for _, a := range strs {
		n = (n << 8) + Dtoi(a)
	}
	return
}

func Dtoi(s string) (n uint32) {
	for i := 0; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			n = n*10 + uint32(s[i]-'0')
		}
	}
	return
}
