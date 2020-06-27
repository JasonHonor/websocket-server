package utils

import (
	"fmt"
	"net"
	"strings"

	"github.com/shirou/gopsutil/host"
)

func GetAdapterInfo() (map[string]string, map[string]string, error) {

	ips := make(map[string]string)
	macs := make(map[string]string)

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			return nil, nil, err
		}
		addresses, err := byName.Addrs()

		var sOldIP string = ""
		for _, v := range addresses {

			sIP := v.String()

			if strings.Contains(sIP, ":") {
				continue
			}

			sOldIP = ips[byName.Name]
			if sOldIP != "" {
				sOldIP += ","
			}

			ips[byName.Name] = sOldIP + sIP
		}

		macAddress := byName.HardwareAddr.String()
		macs[byName.Name] = macAddress
	}
	return ips, macs, nil
}

//GetIPList 获取网卡IP列表
func GetIPList() string {
	nicList, macList, _ := GetAdapterInfo()
	var sRet string
	for name, ip := range nicList {

		if name == "lo" {
			continue
		}

		if sRet != "" {
			sRet += ";"
		}

		sRet += name + "#" + macList[name] + "#" + ip
	}
	return sRet
}

//GetHostName 获取主机名称
func GetSysInfo() string {
	info, _ := host.Info()
	return fmt.Sprintf("%s|%s|%v|%s|%s|%s", info.Hostname, info.OS, info.BootTime, info.Platform, info.PlatformVersion, GetIPList())
}
