package main

import (
	"log"
	"path/filepath"
)

var srcDir = "/mnt/d/miuiback"
var dstDir = "/mnt/d/miuiback2"

func main() {

	dockerCli, err := OpenDockerCli()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("hello")

	// 去除MIUI头文件
	Cut(filepath.Join(srcDir, "wechat.bak"))

	// 解包
	dockerCli.RunWithDelete("greycodee/abe", []string{srcDir + ":/bak"},
		[]string{"java", "-jar", "abe.jar", "unpack", "/bak/wechat_miui.bak", "/bak/unpackGo.tar"})

	err = UnPackTar(filepath.Join(srcDir, "unpackGo.tar"), srcDir)
	if err != nil {
		log.Fatal(err)
	}

	infoList := GetWechatAccountInfo(srcDir)
	ArchiveFile(srcDir, dstDir, infoList, dockerCli)

	// 解压
	// fileprocessor := &FileProcessor{}
	// err := fileprocessor.UnPackTar("/Users/zheng/Documents/miuibak/wechat3_docker.tar", "/Users/zheng/Documents/miuibak/test/")
	// if err != nil {
	// 	log.Println(err)
	// }
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
