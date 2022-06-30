package main

import (
	"archive/tar"
	"bufio"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	VideoParentPath    = "Android/data/com.tencent.mm/MicroMsg/"
	VoiceParentPath    = "Android/data/com.tencent.mm/MicroMsg/"
	ImageParentPath    = "apps/com.tencent.mm/r/MicroMsg/"
	AvatarParentPath   = "apps/com.tencent.mm/r/MicroMsg/"
	DownloadParentPath = "Android/data/com.tencent.mm/MicroMsg/"
	DBParentPath       = "apps/com.tencent.mm/r/MicroMsg/"

	VideoDirName      = "video"
	VoiceDirName      = "voice2"
	ImageDirName      = "image2"
	AvatarDirName     = "avatar"
	DownloadDirName   = "Download"
	EnMicroMsgDBName  = "EnMicroMsg.db"
	WxFileIndexDBName = "WxFileIndex.db"
)

type WechatAccountInfo struct {
	SystemFileUID string
	SDCardFileUID string
	Uin           string
	DBKey         string
}

type UINMap struct {
	Set UINSet `xml:"set"`
}
type UINSet struct {
	Str []string `xml:"string"`
}

type FileProcessor struct {
}

func UnPackTar(src string, dst string) error {

	fr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fr.Close()

	tr := tar.NewReader(fr)

	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}

		dstFileDir := filepath.Join(dst, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if b := ExistDir(dstFileDir); !b {
				if err := os.MkdirAll(dstFileDir, 0775); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			dirname := string([]rune(dstFileDir)[0:strings.LastIndex(dstFileDir, "/")])
			if b := ExistDir(dirname); !b {
				if err := os.MkdirAll(dirname, 0775); err != nil {
					return err
				}
			}
			if err != nil {
				return err
			}
			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}
			fmt.Printf("成功解压: %s\n", dstFileDir)
			file.Close()
		}
	}

}

func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

func CopyFile(src string, dst string) error {
	inStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if inStat.IsDir() {
		err = os.MkdirAll(dst, inStat.Mode())
		if err != nil {
			return err
		}
		fileList, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, v := range fileList {
			CopyFile(filepath.Join(src, v.Name()), filepath.Join(dst, v.Name()))
		}
	} else {
		dstParentPath := string([]rune(dst)[0:strings.LastIndex(dst, "/")])
		err := os.MkdirAll(dstParentPath, 0775)
		if err != nil {
			return err
		}
		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()
		inbuf := bufio.NewReader(in)

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, inbuf)
		if err != nil {
			return err
		}

	}
	return nil
}

func CutMiuiBak(src string) string {
	fin, err := os.Open(src)
	if err != nil {
		log.Fatalln(err)
	}
	defer fin.Close()

	bakHead := make([]byte, 4)
	_, err = fin.Read(bakHead)
	if err != nil {
		log.Fatalln(err)
	}
	if string(bakHead) != "MIUI" {
		log.Fatalln("不是小米的备份文件")
	}

	parentPath := string([]rune(src)[0:strings.LastIndex(src, "/")])
	newFile := filepath.Join(parentPath, "wechat_miui.bak")
	fout, err := os.Create(newFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer fout.Close()

	_, err = fin.Seek(41, io.SeekStart)
	if err != nil {
		log.Fatalln(err)
	}

	n, err := io.Copy(fout, fin)
	log.Printf("Copied %d bytes, err: %v", n, err)
	return newFile
}

func GetWechatAccountInfo(src string) []WechatAccountInfo {
	// 获取uin
	globalSp := filepath.Join(src, "apps/com.tencent.mm/sp/app_brand_global_sp.xml")
	file, err := os.Open(globalSp)
	if err != nil {
		log.Fatalf("open file failed!%s", err)
	}
	defer file.Close()

	uin := make([]byte, 4096)
	n, _ := file.Read(uin)
	log.Println(string(uin[:n]))
	m := UINMap{}
	err = xml.Unmarshal(uin[:n], &m)
	if err != nil {
		log.Println(err)
	}
	accountInfoList := make([]WechatAccountInfo, 0)
	for _, uin := range m.Set.Str {
		info := WechatAccountInfo{}
		info.Uin = uin
		f := []byte("mm" + uin)
		fdName := md5.Sum(f)
		systemFileUID := fmt.Sprintf("%x", fdName)
		log.Println(systemFileUID)
		info.SystemFileUID = systemFileUID
		s := []byte("1234567890ABCDEF" + uin)
		fdName = md5.Sum(s)
		dbKey := fmt.Sprintf("%x", fdName)[:7]
		info.DBKey = dbKey
		log.Printf("DBKEY: %s", dbKey)
		acctMapping := filepath.Join(src, "apps/com.tencent.mm/r/MicroMsg", systemFileUID, "account.mapping")
		am, err := os.Open(acctMapping)
		if err != nil {
			log.Fatalf("open acctMapping file failed!%s", err)
		}
		defer am.Close()

		sdPathBytes := make([]byte, 4096)
		l, _ := am.Read(sdPathBytes)
		sdPathStr := string(sdPathBytes[:l])
		log.Printf("SDCard Path: %s", sdPathStr)
		info.SDCardFileUID = sdPathStr
		accountInfoList = append(accountInfoList, info)
	}
	return accountInfoList
}

func ArchiveFile(src string, dst string, list []WechatAccountInfo, dockerCli DockerCli) error {
	err := os.MkdirAll(dst, 0775)
	if err != nil {
		return err
	}

	for _, v := range list {
		// 移动DB
		enmicromsgFilePath := filepath.Join(dst, v.SystemFileUID)
		CopyFile(filepath.Join(src, DBParentPath, v.SystemFileUID, EnMicroMsgDBName), filepath.Join(enmicromsgFilePath, EnMicroMsgDBName))
		wxfileindexFilePath := filepath.Join(dst, v.SystemFileUID)
		CopyFile(filepath.Join(src, DBParentPath, v.SystemFileUID, WxFileIndexDBName), filepath.Join(wxfileindexFilePath, WxFileIndexDBName))
		// 移动 Video
		CopyFile(filepath.Join(src, VideoParentPath, v.SDCardFileUID, VideoDirName), filepath.Join(dst, v.SystemFileUID, VideoDirName))
		// 移动 voice
		voiceFilePath := filepath.Join(dst, v.SystemFileUID, VoiceDirName)
		CopyFile(filepath.Join(src, VoiceParentPath, v.SDCardFileUID, VoiceDirName), voiceFilePath)
		// 移动 avatar
		CopyFile(filepath.Join(src, AvatarParentPath, v.SystemFileUID, AvatarDirName), filepath.Join(dst, v.SystemFileUID, AvatarDirName))
		// 移动images
		CopyFile(filepath.Join(src, ImageParentPath, v.SystemFileUID, ImageDirName), filepath.Join(dst, v.SystemFileUID, ImageDirName))
		// 移动 Download文件夹
		CopyFile(filepath.Join(src, DownloadParentPath, DownloadDirName), filepath.Join(dst, v.SystemFileUID, DownloadDirName))

		dockerCli.RunWithDelete("greycodee/wcdb-sqlcipher", []string{enmicromsgFilePath + ":" + "/wcdb"}, []string{"-f", EnMicroMsgDBName, "-k", v.DBKey})
		dockerCli.RunWithDelete("greycodee/wcdb-sqlcipher", []string{wxfileindexFilePath + ":" + "/wcdb"}, []string{"-f", WxFileIndexDBName, "-k", v.DBKey})
		dockerCli.RunWithDelete("greycodee/silkv3-decoder", []string{voiceFilePath + ":/media"}, nil)
	}
	return nil
}
