package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type UINMap struct {
	Set UINSet `xml:"set"`
}
type UINSet struct {
	Str []string `xml:"string"`
}

func main() {
	// fmt.Println("hello")

	// 去除MIUI头文件
	// cut()

	// 解压
	fileprocessor := &FileProcessor{}
	err := fileprocessor.UnPackTar("/Users/zheng/Documents/miuibak/wechat3_docker.tar", "/Users/zheng/Documents/miuibak/test/")
	if err != nil {
		log.Println(err)
	}
	// cmd := exec.Command("tar", "-zxvf", "/Users/zheng/Documents/miuibak/wechat3_docker.tar", "-C", "/Users/zheng/Documents/miuibak")
	// cmd.Run()

	// 获取uin
	// file, err := os.Open("/Users/zheng/Documents/miuibak/apps/com.tencent.mm/sp/app_brand_global_sp.xml")
	// if err != nil {
	// 	log.Printf("open file failed!%s", err)
	// }
	// uin := make([]byte, 2048)
	// n, _ := file.Read(uin)
	// log.Println(string(uin[:n]))
	// m := UINMap{}
	// err = xml.Unmarshal(uin[:n], &m)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(m.Set.Str[0])
	// // 获取所有文件夹名称
	// for _, uin := range m.Set.Str {
	// 	f := []byte("mm" + uin)
	// 	fdName := md5.Sum(f)
	// 	log.Println(fmt.Sprintf("%x", fdName))
	// 	s := []byte("1234567890ABCDEF" + uin)
	// 	fdName = md5.Sum(s)
	// 	log.Println(fmt.Sprintf("%x", fdName)[:7])
	// }

	// 转移文件

	// cli, err := client.NewClientWithOpts(client.FromEnv)
	// if err != nil {
	// 	panic(err)
	// }

	// dockerRUN(*cli, "greycodee/wcdb-sqlcipher", []string{"/Users/zheng/coding/study/miui-wechat-bak/apps/com.tencent.mm/r/MicroMsg/79b23ef49a3016d8c52a787fc4ab59e4:/wcdb"},
	// 	[]string{"-f", "EnMicroMsg.db", "-k", "626d0bc"})

	// dockerRUN(*cli, "greycodee/abe:1.0", []string{"/Users/zheng/Documents/miuibak/:/bak"},
	// 	[]string{"java", "-jar", "abe.jar", "unpack", "/bak/wechat_new.bak", "/bak/unpackGo2.tar"})

	// err = cli.ContainerStart(context.Background(), "b88a266e041661e2890b43a978dac25068c36f0de82024ba8e3ebc9bb3d5c8b6", types.ContainerStartOptions{})
	// if err != nil {
	// 	log.Println(err)
	// }

}

func dockerRUN(cli client.Client, imageName string, bind []string, cmd []string) {

	// 拉取镜像
	// log.Println("正在拉取 Docker 镜像...")
	// read, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	// if err != nil {
	// 	log.Fatalf("拉取镜像失败: %s\n", err)
	// } else {
	// 	log.Print("拉取镜像成功:")
	// 	pullLog, _ := ioutil.ReadAll(read)
	// 	log.Println(string(pullLog))
	// }

	// 创建容器
	log.Println("正在创建容器...")
	containerBody, err := cli.ContainerCreate(context.Background(), &container.Config{
		Cmd:   cmd,
		Image: imageName,
	}, &container.HostConfig{
		Binds: bind,
	}, &network.NetworkingConfig{}, &v1.Platform{}, "")
	if err != nil {
		log.Fatalf("创建容器失败:%s", err)
	} else {
		log.Printf("容器创建成功: %s", containerBody.ID)
	}

	// 启动容器
	log.Println("正在启动容器...")
	err = cli.ContainerStart(context.Background(), containerBody.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatalf("启动容器失败:%s", err)
	} else {
		log.Println("容器启动成功")
	}
	// 容器日志输出
	reader, err := cli.ContainerLogs(context.Background(), containerBody.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Details: true})
	if err != nil {
		log.Fatalf("容器日志输出失败: %s", err)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil && err != io.EOF {
		log.Fatalf("读取失败: %s", err)
	}
	// 删除容器
	log.Println("正在删除容器...")
	log.Printf("容器ID: %s", containerBody.ID)
	err = cli.ContainerRemove(context.Background(), containerBody.ID, types.ContainerRemoveOptions{})
	if err != nil {
		log.Fatalf("删除容器失败: %s", err)
	} else {
		log.Println("删除容器成功")
	}
	// docker run --rm  -v /Users/zheng/Documents/miuibak/:/bak greycodee/abe:1.0 java -Xms1G -Xmx1G -jar abe.jar unpack /bak/wechat_new.bak /bak/wechat3_docker.tar

}

func cut() {
	fin, err := os.Open("/Users/zheng/Documents/miuibak/wechat.bak")
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	fout, err := os.Create("/Users/zheng/Documents/miuibak/wechat_new.bak")
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	// Offset is the number of bytes you want to exclude
	_, err = fin.Seek(41, io.SeekStart)
	if err != nil {
		panic(err)
	}

	n, err := io.Copy(fout, fin)
	fmt.Printf("Copied %d bytes, err: %v", n, err)
}
