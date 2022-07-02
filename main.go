package main

import (
	"flag"
	"log"
	"path/filepath"
)

var srcDir = flag.String("s", "", "src dir")
var dstDir = flag.String("d", "", "dst dir")

func init() {
	flag.Parse()
	if *srcDir == "" || *dstDir == "" {
		log.Fatalln("请输入src和dst目录")
	}
}

func main() {

	dockerCli, err := OpenDockerCli()
	if err != nil {
		log.Fatal(err)
	}
	Unzip(filepath.Join(*srcDir, "backup_wechat.zip"), *srcDir)

	// 去除MIUI头文件
	CutMiuiBak(filepath.Join(*srcDir, "wechat.bak"))

	// 解包
	dockerCli.RunWithDelete("greycodee/abe", []string{*srcDir + ":/bak"},
		[]string{"java", "-jar", "abe.jar", "unpack", "/bak/wechat_miui.bak", "/bak/unpackGo.tar"})

	err = UnPackTar(filepath.Join(*srcDir, "unpackGo.tar"), *srcDir)
	if err != nil {
		log.Fatal(err)
	}

	infoList := GetWechatAccountInfo(*srcDir)
	ArchiveFile(*srcDir, *dstDir, infoList, dockerCli)
}
