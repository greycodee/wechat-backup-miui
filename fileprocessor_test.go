package main

import (
	"fmt"
	"log"
	"testing"
)

func TestCopyFile(t *testing.T) {
	err := CopyFile("/Users/zheng/coding/study/wechat-backup-miui/test", "/Users/zheng/coding/study/wechat-backup-miui/test2")
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetWechatAccountInfo(t *testing.T) {
	infoList := GetWechatAccountInfo("/Users/zheng/Documents/miuibak")
	for _, v := range infoList {
		fmt.Println(v)
	}
}

func TestArchiveFile(t *testing.T) {

	infoList := GetWechatAccountInfo("/Users/zheng/Documents/miuibak")
	// for _, v := range infoList {

	// }
	dockerCli, err := OpenDockerCli()
	if err != nil {
		log.Fatal(err)
	}
	ArchiveFile("/Users/zheng/Documents/miuibak", "/Users/zheng/Documents/miuibak3", infoList, dockerCli)
}

func TestUnzip(t *testing.T) {
	Unzip("/Users/zheng/Documents/miuibak/backup_wechat.zip", "/Users/zheng/Documents/miuibak")
}

func TestCutMiuiBak(t *testing.T) {
	CutMiuiBak("/Users/zheng/Documents/miuibak/wechat.bak")
}
