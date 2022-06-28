package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileProcessor struct {
}

func (f FileProcessor) UnPackTar(src string, dst string) error {

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
