// @Author: YangPing
// @Create: 2023/10/23
// @Description: IP工具类

package util

import (
	"fmt"
	"net"
	"strings"
)

const (
	IPV4    = "IPV4"
	IPV6    = "IPV6"
	IPV4SEG = "IPV4_SEG"
	IPV6SEG = "IPV6_SEG"
)

func IpType(ipStr string) string {
	if ipStr == "" {
		return ipStr
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		_, ipNet, err := net.ParseCIDR(ipStr)
		if err == nil {
			if ipNet.IP.To4() != nil {
				return IPV4SEG
			} else {
				return IPV6SEG
			}
		}
		return ""
	}
	to4 := ip.To4()
	if to4 != nil {
		return IPV4
	}
	to6 := ip.To16()
	if to6 != nil {
		return IPV6
	}
	return ""
}

func IPTo16Str(address string) string {
	if address == "" {
		return address
	}
	ip := net.ParseIP(address)
	if ip == nil {
		return ""
	}
	to4 := ip.To4()
	if to4 != nil {
		return CompleteIPv4(ip)
	}
	to6 := ip.To16()
	if to6 != nil {
		return CompleteIPv6(ip)
	}
	return address
}

func IpCheck(address string) bool {
	ip := net.ParseIP(address)
	if ip != nil {
		return true
	}
	_, _, err := net.ParseCIDR(address)
	if err != nil {
		return false
	}
	return true
}

func GetIPRange(ipRange string) (string, string) {
	ip, ipNet, err := net.ParseCIDR(ipRange)
	if err != nil {
		return "", ""
	}
	// 获取开始IP地址
	startIP := ip.Mask(ipNet.Mask)
	// 获取结束IP地址
	endIP := make(net.IP, len(startIP))
	copy(endIP, startIP)
	for i := len(endIP) - 1; i >= 0; i-- {
		if ipNet.Mask[i] == 0xff {
			break
		}
		endIP[i] |= ^ipNet.Mask[i]
	}
	to4 := startIP.To4()
	if to4 != nil {
		return CompleteIPv4(startIP), CompleteIPv4(endIP)
	} else {
		return CompleteIPv6(startIP), CompleteIPv6(endIP)
	}

}

func CompleteIPv4(parsedIP net.IP) string {
	return fmt.Sprintf("%02x%02x%02x%02x", parsedIP.To4()[0], parsedIP.To4()[1], parsedIP.To4()[2], parsedIP.To4()[3])
}

func CompleteIPv6(ip net.IP) string {
	ipStrComplete := ip.String()
	if strings.Contains(ipStrComplete, "::") {
		segments := strings.Split(ipStrComplete, "::")
		segmentsCount := 8 - len(strings.Split(segments[0], ":")) - len(strings.Split(segments[1], ":"))
		ipStrComplete = strings.Replace(ipStrComplete, "::", ":"+strings.Repeat("0000:", segmentsCount), 1)
	}
	if len(ipStrComplete) == 39 {
		// 完整格式
		return strings.ReplaceAll(ipStrComplete, ":", "")
	}
	var str []string
	for _, v := range strings.Split(ipStrComplete, ":") {
		v = fmt.Sprintf("%04s", v)
		str = append(str, v)
	}
	return strings.Join(str, "")
}
