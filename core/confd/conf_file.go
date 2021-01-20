package confd

/**
 * Copyright 2021  gowrk Author. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @Project gowk
 * @Description go framework
 * @author XiongChuanLiang<br/>(xcl_168@aliyun.com)
 * @license http://www.apache.org/licenses/  Apache v2 License
 * @version 1.0
 */

import (
	"fmt"
	"io"
	"os"

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
	path, ext, err := PaseConfigFile(confFile)
	if err != nil {
		return err
	}

	if loadViper == nil {
		return fmt.Errorf(" viper is null.")
	}

	loadViper.AddConfigPath(path)
	path2, _ := os.Getwd()
	loadViper.AddConfigPath(path2)
	loadViper.SetConfigType(cf.getConfigType(ext))
	//读取配置
	err = loadViper.ReadInConfig()
	if err != nil {
		return err
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
		return err
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
