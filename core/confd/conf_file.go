package confd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// //Config -
type ConfLocalFile struct {
}

func NewConfLocalFile() *ConfLocalFile {
	return &ConfLocalFile{}
}

//
func (cf *ConfLocalFile) LoadConfigFromLocalFile(loadViper *viper.Viper, confFile string) error {
	if loadViper == nil {
		return fmt.Errorf(" viper is null.")
	}

	path, ext, err := PaseConfigFile(confFile)
	if err != nil {
		return fmt.Errorf(" err:%s confFile:%s",err,confFile)
	}
	if path != "" {
		loadViper.AddConfigPath(path)
	}
	path2, _ := os.Getwd()
	loadViper.AddConfigPath(path2)
	loadViper.AddConfigPath(".")
	loadViper.SetConfigName(filepath.Base(confFile))
	loadViper.SetConfigType(cf.getConfigType(ext))
	//读取配置
	err = loadViper.ReadInConfig()
	if err != nil {
		return fmt.Errorf(" viper 读取本地配置文件失败. err:%s path:%s", err.Error(),path2)
	}

	return nil
}

func (cf *ConfLocalFile) ReadConfig(loadViper *viper.Viper, in io.Reader, configType string) error {

	if loadViper == nil {
		return fmt.Errorf(" viper is null.")
	}

	loadViper.SetConfigType(cf.getConfigType(configType))
	err := loadViper.ReadConfig(in)
	if err != nil {
		return fmt.Errorf(" viper 读取失败. err:%s configType:%s", err.Error(), configType)
	}

	return nil
}

func (cf *ConfLocalFile) getConfigType(configType string) string {
	switch configType {
	case "json", "hcl", "prop", "props", "properties", "dotenv", "env", "toml", "yaml", "yml", "ini":
		return configType
	default:
		return _ConfigType
	}
}

//PaseConfigFile 解析文件，得到path,ext
func PaseConfigFile(confFile string) (string, string, error) {
	// if _, err := os.Stat(confFile); os.IsNotExist(err) {
	// 	return "", "", err
	// }
	path, fileName := filepath.Split(confFile)
	ext := strings.ToLower(filepath.Ext(fileName))
	fext := SubString(ext, 1, len(ext))

	return path, fext, nil
}

//包含中文的字符串截取
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
