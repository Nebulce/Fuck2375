package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// 计算192.168.1.1/24类的IP地址
func getAllIPsInCIDR(ipNet *net.IPNet)  {

	// 遍历 CIDR 地址范围内的所有 IP 地址
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc1(ip) {
		IP = append(IP,ip.String())
	}

}

// 增加 192.168.1.1/24类的IP地址
func inc1(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

//计算192.168.1.1-255类的IP地址
func parseIPRange(IPRange string){
	// 定义正则表达式
	re := regexp.MustCompile(`(\d+)-(\d+)`)

	// 使用正则表达式提取匹配的子串
	matches := re.FindStringSubmatch(IPRange)

	// 定义正则表达式
	re1 := regexp.MustCompile(`^(\d+\.\d+\.\d+)\.\d+-\d+$`)

	// 使用正则表达式提取匹配的子串
	matches1 := re1.FindStringSubmatch(IPRange)

	// 输出匹配的结果
	if len(matches) == 3 {
		start,_ := strconv.Atoi(matches[1])
		end,_ := strconv.Atoi(matches[2])
		for i := start; i <= end; i++ {
			tmp :=(matches1[1]  + "." + strconv.Itoa(i))
			IP = append(IP,tmp)
		}
	} else {
		fmt.Println("未找到匹配的模式")
	}
}

//获取用户输入的IP地址
func getIP(ip string){
	substrings := strings.Split(ip, ",")

	// 将分割出的每个字符串添加到切片中
	for _, substring := range substrings {
		ipPattern := regexp.MustCompile(`\b(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)
		if(ipPattern.FindAllString(substring, -1))!=nil{
			if strings.Contains(substring, "/") {

				_, ipNet, err := net.ParseCIDR(substring)
				if err != nil {
					fmt.Println("无法解析 CIDR 地址:", err)
					return
				}
				getAllIPsInCIDR(ipNet)
			}else if strings.Contains(substring, "-") {
				parseIPRange(substring)
			}else{
				IP = append(IP,substring)
			}
		}
	}
}
