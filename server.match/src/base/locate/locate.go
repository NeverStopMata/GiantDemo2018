package locate

import (
	"base/glog"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	nettypes = []string{"未知", "电信", "联通", "铁通", "移动", "教育网", "长城宽带"}
)

type Location struct {
	Sip uint32 // 起始ip
	Eip uint32 // 结束ip
	Cid uint32 // 国家城市编号
	Net uint8  // 网络类型
}

func Inet_ntoa(ipnr uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte((ipnr>>24)&0xFF), byte((ipnr>>16)&0xFF), byte((ipnr>>8)&0xFF), byte(ipnr&0xFF))
}

func Inet_aton(ipaddr string) uint32 {
	bits := strings.Split(ipaddr, ".")
	if len(bits) != 4 {
		glog.Error("[位置] 格式错误 ", ipaddr)
		return 0
	}
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	return uint32(b0)<<24 | uint32(b1)<<16 | uint32(b2)<<8 | uint32(b3)
}

var ips []Location

func Load(fname string) bool {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		glog.Error("[位置] 打开配置失败 ", fname, ",", err)
		return false
	}

	buffer := bytes.NewBuffer(content)

	for {
		ldata := Location{}
		err := binary.Read(buffer, binary.BigEndian, &ldata)
		if err != nil {
			break
		}
		ips = append(ips, ldata)
	}

	//for k, v := range ips {
	//	fmt.Println(k, " ", inet_ntoa(v.Sip), ",", v.Sip, " ", inet_ntoa(v.Eip), ",", v.Eip, " ", v.Cid>>16, ",", v.Cid&0x0000ffff)
	//}

	glog.Info("[位置] 读取配置成功 ", fname, ",", len(ips))
	return true
}

func GetLoc(sip string) (uint32, uint8) {
	if len(ips) == 0 {
		return 0, 0
	}

	ip := Inet_aton(sip)
	if ip == 0 {
		return 0, 0
	}

	var start int = 0
	var end int = len(ips) - 1

	for start <= end {
		mid := (start + end) / 2
		if mid >= len(ips) {
			return 0, 0
		}
		if ips[mid].Sip > ip {
			end = mid - 1
		} else if ips[mid].Eip < ip {
			start = mid + 1
		} else {
			tip := ips[mid]
			return tip.Cid, tip.Net
		}
	}

	return 0, 0
}

func GetLocCity(sip string) (uint32, uint32, uint8) {
	if len(ips) == 0 {
		return 0, 0, 0
	}

	ip := Inet_aton(sip)
	if ip == 0 {
		return 0, 0, 0
	}

	var start int = 0
	var end int = len(ips) - 1

	for start <= end {
		mid := (start + end) / 2
		if mid >= len(ips) {
			return 0, 0, 0
		}
		if ips[mid].Sip > ip {
			end = mid - 1
		} else if ips[mid].Eip < ip {
			start = mid + 1
		} else {
			tip := ips[mid]
			return uint32(tip.Cid>>16), uint32(tip.Cid&0x0000ffff), tip.Net
		}
	}

	return 0, 0, 0
}


func PToCidr(text string) (uint32, uint32) {
	addrs := strings.Split(text, "/")
	if len(addrs) != 2 {
		return Inet_aton(text), 0xffffffff
	}
	shift, _ := strconv.Atoi(addrs[1])
	if shift > 32 {
		return Inet_aton(addrs[0]), 0xffffffff
	}
	return Inet_aton(addrs[0]), uint32(0xffffffff << (32 - uint32(shift)))
}

// 是否是电信
func IsTele(cnet uint8) bool {
	return cnet == 4
}

type CityCode struct {
	City string `xml:"Name,attr"`
	Code uint32 `xml:"Id,attr"`
}

type CityCodeXml struct {
	Items []CityCode `xml:"item"`
}

var CityCodeMap map[string]uint32

func InitCityCode(cfg string) bool {
    CityCodeMap = make(map[string]uint32)
	content, err := ioutil.ReadFile(cfg)
	if err != nil {
		glog.Error("[配置] 打开地区配置失败 ", err)
		return false
	}
	var codeXml CityCodeXml
	err = xml.Unmarshal(content, &codeXml)
	if err != nil {
		glog.Error("[配置]解析地区配置失败:", err)
		return false
	}
	for _, code := range codeXml.Items {
		CityCodeMap[code.City] = code.Code
	}
	return true
}

func GetCityCode(city string) (uint32, bool) {
	cCode := CityCodeMap[city]
	if cCode == 0 {
		return 0, false
	}
	return 10223616 + cCode, true
}
