package main

import (
	"archive/zip"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerCli struct {
	cli client.Client
}

func OpenDockerCli() (DockerCli, error) {
	dockerCli := &DockerCli{}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return *dockerCli, err
	}
	dockerCli.cli = *cli
	return *dockerCli, nil
}

func (cli DockerCli) RunWithDelete(imageName string, bind []string, cmd []string) {
	// 拉取镜像
	log.Println("正在拉取 Docker 镜像...")
	read, err := cli.cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		log.Fatalf("拉取镜像失败: %s\n", err)
	} else {
		log.Print("拉取镜像成功:")
		pullLog, _ := ioutil.ReadAll(read)
		log.Println(string(pullLog))
	}

	// 创建容器
	log.Println("正在创建容器...")
	containerBody, err := cli.cli.ContainerCreate(context.Background(), &container.Config{
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
	err = cli.cli.ContainerStart(context.Background(), containerBody.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatalf("启动容器失败:%s", err)
	} else {
		log.Println("容器启动成功")
	}
	// 容器日志输出
	reader, err := cli.cli.ContainerLogs(context.Background(), containerBody.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Details: true})
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
	err = cli.cli.ContainerRemove(context.Background(), containerBody.ID, types.ContainerRemoveOptions{})
	if err != nil {
		log.Fatalf("删除容器失败: %s", err)
	} else {
		log.Println("删除容器成功")
	}
}

func Unzip(zipPath, dstDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if err := unzipFile(file, dstDir); err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(file *zip.File, dstDir string) error {
	filePath := path.Join(dstDir, file.Name)
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, rc)
	return err
}
