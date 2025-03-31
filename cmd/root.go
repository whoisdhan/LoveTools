package cmd

import (
	"github.com/spf13/cobra"
)

var (
	targets  []string                  // 目标列表，可以是域名或 IP
	proxy    string                    // 代理地址
	stopChan chan bool                 // 停止加载动画的通道
	dict     string                    // 字典文件
	yamlPath string    = "config.yaml" // yaml配置文件
)

var rootCmd = &cobra.Command{
	Use: "LoveTools",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help() // 显示帮助信息
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&targets, "target", "T", nil, "目标域名或IP")              // 添加目标参数
	rootCmd.PersistentFlags().StringVarP(&proxy, "proxy", "P", "", "e.g. http(s)://127.0.0.1:8080") // 添加代理参数
	rootCmd.PersistentFlags().StringVarP(&dict, "dict", "F", "dict.txt", "字典文件,默认为dict.txt")        // 添加字典参数
}
