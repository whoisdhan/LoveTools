package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DomainLocalDict string              `yaml:"DomainLocalDict"` // 本地域名字典路径
	UrlDict         string              `yaml:"UrlDict"`         // 域名字典url
	CDNList         map[string][]string `yaml:"CDNList"`
}

func newConfig() *Config {
	return &Config{}
}

// 解析yaml，返回Config结构体指针
func ParseConfig(path string) *Config {
	var err error
	var fileinfo os.FileInfo
	fileinfo, err = os.Stat(path)
	if err != nil {
		fmt.Println("找不到文件")
	}
	if fileinfo.IsDir() {
		fmt.Println("输入的路径是文件夹")
	}
	yamlFile, _ := os.ReadFile(path)

	config := newConfig()
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		fmt.Println("yaml文件解析失败")
	}
	return config
}
