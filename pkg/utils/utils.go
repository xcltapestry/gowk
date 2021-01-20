package utils

import (
	"os"
	"path/filepath"
	"strings"
)

//PaseConfigFile 解析文件，得到path,ext
func PaseConfigFile(confFile string) (string, string, error) {
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		return "", "", err
	}
	path, fileName := filepath.Split(confFile)
	ext := strings.ToLower(filepath.Ext(fileName))
	fext := SubString(ext, 1, len(ext))

	return path, fext, nil
}

//IsNotExist 文件是否存在
func IsNotExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return true
	}
	return false
}

//SubString 能处理包含中文的字符串截取
func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}
