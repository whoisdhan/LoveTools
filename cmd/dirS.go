package cmd

import (
	"LoveTools/util"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
)

var (
	dirTimeOut int // 超时时间
)

// 判断是否存在http，不存在就默认加上http
func checkHttp(target string) string {
	if !strings.HasPrefix(target, "http") {
		return "http://" + target
	}
	return target
}

func dirScan(targets []string, dict string) {
	stopChan = make(chan bool)    // 创建一个通道用于停止加载动画
	go util.ShowLoading(stopChan) // 启动加载动画的协程
	client := req.C().
		SetRedirectPolicy(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}).                                                 // 禁止自动跳转
		SetTimeout(time.Duration(dirTimeOut) * time.Second) // 设置超时时间

	if proxy != "" {
		client.SetProxyURL(proxy) // 设置代理
	}

	file, err := os.Open(dict) // 打开字典文件
	defer file.Close()         // 延迟关闭文件
	if err != nil {
		fmt.Println("打开字典文件失败:", err)
		return
	}

	for _, target := range targets {
		target = checkHttp(target) // 检查是否以http开头，如果没有则添加

		//按行读取字典文件
		scanner := bufio.NewScanner(file)

		fmt.Printf("\n正在扫描 %s\n", target)
		for scanner.Scan() {
			line := scanner.Text() // 读取一行
			if line == "" {
				continue // 如果是空行，跳过
			}
			line = util.TrimSlashProper(strings.TrimSpace(line)) // 去除首尾空格
			// 拼接URL
			url := fmt.Sprintf("%s/%s", target, line)
			res, err := client.R().Get(url) // 发送请求
			if err != nil {
				panic(err)
			}

			//根据状态码输出不同颜色的结果
			var path string // 定义路径变量
			if res.IsSuccessState() {
				path = fmt.Sprintf("\r\033[K%s   %d\n", url, res.StatusCode)
				path = color.GreenString(path) // 成功状态码，绿色输出
			}
			if res.StatusCode >= 300 && res.StatusCode < 400 {
				path = fmt.Sprintf("\r\033[K%s   %d\n", url, res.StatusCode)
				path = color.BlueString(path) // 重定向状态码，蓝色输出
			}
			if res.StatusCode >= 400 && res.StatusCode < 500 {
				path := fmt.Sprintf("\r\033[K%s   %d\n", url, res.StatusCode)
				path = color.YellowString(path) // 客户端错误，黄色输出
			}
			if res.StatusCode >= 500 && res.StatusCode < 600 {
				path := fmt.Sprintf("\r\033[K%s   %d\n", url, res.StatusCode)
				path = color.RedString(path) // 服务器错误，红色输出
			}
			fmt.Print(path) // 打印路径
		}
		file.Seek(0, 0) // 重置文件指针到开头
		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err) // 处理扫描错误
		}
	}
	stopChan <- true // 停止加载动画
}

var dirS = &cobra.Command{
	Use:   "dir",
	Short: "目录扫描",
	Long:  `目录扫描,自定义字典`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(targets) == 0 {
			cmd.Help() // 显示帮助信息
			return
		}
		dirScan(targets, dict) // 执行目录扫描
	},
}

func init() {
	rootCmd.AddCommand(dirS)                                     // 将 dirS 命令添加到根命令中
	dirS.Flags().IntVarP(&dirTimeOut, "timeout", "m", 5, "超时时间") // 设置超时时间
}
