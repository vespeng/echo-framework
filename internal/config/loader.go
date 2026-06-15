package config

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	vortego "echo-framework"

	"github.com/spf13/viper"
)

var (
	conf    *Config
	once    sync.Once
	confErr error
)

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	once.Do(func() {
		var err error

		// 配置文件加载路径规则（优先级从高到低）：
		// 1. 外部配置文件：通过明确指定的完整路径加载（默认项目同级目录下config.yaml）
		// 2. 内部配置文件：嵌入到程序中的默认配置文件。
		exConfigFilePath := "../config.yaml"

		viper.SetConfigType("yaml")

		if isConfigFileExist(exConfigFilePath) {
			viper.SetConfigFile(exConfigFilePath)
			err = viper.ReadInConfig()
			if err != nil {
				confErr = fmt.Errorf("failed to read the configuration file: %w", err)
				return
			}
		} else {
			configFile, cErr := vortego.ConfigFile.ReadFile("internal/config/config.default.yaml")
			if cErr != nil {
				confErr = fmt.Errorf("failed to read embedded config file: %w", cErr)
				return
			}

			err = viper.ReadConfig(bytes.NewReader(configFile))
			if err != nil {
				confErr = fmt.Errorf("failed to parse the configuration file: %w", err)
				return
			}
		}

		conf = &Config{}
		err = viper.Unmarshal(conf)
		if err != nil {
			confErr = fmt.Errorf("failed to parse the configuration file: %w", err)
			return
		}
		// 打印系统版本号
		fmt.Printf("Vortego version:  %s\n", conf.App.Version)
	})
	return conf, confErr
}

// isConfigFileExist 检查配置文件是否存在
func isConfigFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
