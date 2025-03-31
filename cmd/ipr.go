package cmd

import (
	"LoveTools/util"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
)

// 反查IP地址调用函数
func ipr() {
	stopChan = make(chan bool)    // 创建一个通道用于停止加载动画
	go util.ShowLoading(stopChan) // 启动加载动画的协程
	var res [][]string
	res = iprsSearch(targets) // 反查IP地址
	stopChan <- true          // 停止加载动画
	util.PrettyPrint(res)     // 打印结果
}

func iprSearch(target string, client *req.Client) []string {

	//反查结果
	var (
		domain  string
		address string
	)

	data, _ := client.R().
		SetPathParam("target", target).
		Get("https://site.ip138.com/{target}")
	doc, _ := goquery.NewDocumentFromReader(data.Body)
	doc.Find("#list > li:nth-child(3) > a").Each(func(i int, s *goquery.Selection) { //#list > li:nth-child(3) > a
		domain = strings.TrimSpace(s.Text()) //域名结果赋值
	})
	doc.Find("h3").Each(func(i int, s *goquery.Selection) { //#list > li:nth-child(3) > a
		address = strings.TrimSpace(s.Text()) //域名结果赋值
	})

	return []string{target, domain, address}
}
func iprsSearch(targets []string) [][]string {
	client := req.C().
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	if proxy != "" {
		client.SetProxyURL(proxy) // 设置代理
	}
	res := [][]string{
		{"目标", "反查结果", "归属地"}, // 表头
	}
	for _, target := range targets {
		tmp := iprSearch(target, client) // 反查IP地址
		res = append(res, tmp)           // 将结果添加到结果列表中
	}
	return res
}

var iprCmd = &cobra.Command{
	Use:   "ipr",
	Short: "IP反查",
	Long:  `IP反查，支持域名和 IP 查询`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(targets) == 0 {
			fmt.Println("缺少参数，请使用 -t 或 --target 指定目标")
			return
		}
		ipr()
	},
}

func init() {
	rootCmd.AddCommand(iprCmd) // 将 iprCmd 添加到根命令中
}
