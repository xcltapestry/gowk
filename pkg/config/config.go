package config

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
	"io"
	"time"

	"github.com/spf13/viper"
)

const _defaultConfigType = "yaml"

//Config -
type Config struct {
	viper *viper.Viper
}

func New() *Config {
	cfg := &Config{}
	if cfg.viper == nil {
		cfg.viper = viper.New()
		cfg.viper.SetConfigType(_defaultConfigType)
	}
	return cfg
}

func (cfg *Config) AddConfigPath(path string) {
	if cfg.viper == nil {
		cfg.viper = viper.New()
		cfg.viper.SetConfigType(_defaultConfigType)
	}
	cfg.viper.AddConfigPath(path)
}

func (cfg *Config) SetConfigType(configType string) {
	if cfg.viper == nil {
		cfg.viper = viper.New()
	}

	switch configType {
	case "json", "hcl", "prop", "props", "properties", "dotenv", "env", "toml", "yaml", "yml", "ini":
		cfg.viper.SetConfigType(configType)
	default:
		cfg.viper.SetConfigType(_defaultConfigType)
	}
}

func (cfg *Config) ReadInConfig(path, name string) error {
	if cfg.viper == nil {
		cfg.viper = viper.New()
		cfg.viper.SetConfigType(_defaultConfigType)
	}
	cfg.viper.SetConfigName(name)
	cfg.viper.AddConfigPath(path)

	// if err := cfg.viper.ReadInConfig(); err != nil {
	// 	// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
	// 	// }
	// 	return err
	// }

	return cfg.viper.ReadInConfig()
}

func (cfg *Config) ReadConfig(in io.Reader) error {
	if cfg.viper == nil {
		cfg.viper = viper.New()
		cfg.viper.SetConfigType(_defaultConfigType)
	}
	return cfg.viper.ReadConfig(in)
}

func (cfg *Config) Get(key string) interface{} {
	if cfg.viper == nil {
		return nil
	}
	return cfg.viper.Get(key)
}

func (cfg *Config) GetString(key string) string {
	if cfg.viper == nil {
		return ""
	}
	return cfg.viper.GetString(key)
}

func (cfg *Config) GetBool(key string) bool {
	return cfg.viper.GetBool(key)
}

func (cfg *Config) GetDuration(key string) time.Duration {
	return cfg.viper.GetDuration(key)
}

func (cfg *Config) GetInt(key string) int {
	return cfg.viper.GetInt(key)
}

func (cfg *Config) GetStringMap(key string) map[string]interface{} {
	return cfg.viper.GetStringMap(key)
}

func (cfg *Config) GetStringMapString(key string) map[string]string {
	return cfg.viper.GetStringMapString(key)
}

func (cfg *Config) GetStringMapStringSlice(key string) map[string][]string {
	return cfg.viper.GetStringMapStringSlice(key)
}

func (cfg *Config) GetStringSlice(key string) []string {
	return cfg.viper.GetStringSlice(key)
}

func (cfg *Config) InConfig(key string) bool {
	return cfg.viper.InConfig(key)
}

func (cfg *Config) IsSet(key string) bool {
	return cfg.viper.IsSet(key)
}
