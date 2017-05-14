package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
)

// 计算文件MD%
func Md5SumFile(file string) (value [md5.Size]byte, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	value = md5.Sum(data)
	return
}

// 文件目录是否存在
func IsDirExist(path string) bool {
	p, err := os.Stat(path)
	fmt.Println(p)
	if err != nil {
		return os.IsExist(err)
	} else {
		return true
	}
}
