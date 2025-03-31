package cmd

import (
	"LoveTools/util"
	"io"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
)

type cdnInfo struct {
	domain      string
	address     []string
	cdnURLList  []string // CDN URL 列表
	reponseList []string // 响应结果列表
	//CDN URL 和 响应结果 个数是一样的
	ok []string // 检测结果
}

var cdnInfoList []cdnInfo

var cdnCmd = &cobra.Command{
	Use:   "cdn",
	Short: "CDN检测",
	Long:  `CDN检测`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(targets) == 0 {
			cmd.Help() // 显示帮助信息
			return
		}
		cdnInfoList = cdns(targets)
		printCDNInfos(cdnInfoList)
		// cdns(targets) // 执行CDN检测
	},
}

func cdnCheck(target string, client *req.Client) cdnInfo {
	var cdnInfo cdnInfo
	config := util.ParseConfig(yamlPath) // 读取配置文件
	for address, cdnURLS := range config.CDNList {
		for _, cdnURL := range cdnURLS {
			req, _ := client.R().
				SetQueryParam("domain", target).
				Get(cdnURL)
			if req.StatusCode == 200 {
				contentBytes, _ := io.ReadAll(req.Body)
				resultIP := string(contentBytes)
				cdnInfo.cdnURLList = append(cdnInfo.cdnURLList, cdnURL)     //将文件解析出来的cdn也放进去结构体，后面要打印
				cdnInfo.reponseList = append(cdnInfo.reponseList, resultIP) //将一个cdn查询结果传进去
				cdnInfo.ok = append(cdnInfo.ok, "检测成功")
			} else {
				cdnInfo.cdnURLList = append(cdnInfo.cdnURLList, cdnURL)  //将文件解析出来的cdn也放进去结构体，后面要打印
				cdnInfo.reponseList = append(cdnInfo.reponseList, "---") //将一个cdn查询结果传进去
				cdnInfo.ok = append(cdnInfo.ok, "检测失败")
			}
			cdnInfo.domain = target
			cdnInfo.address = append(cdnInfo.address, address)
		}
	}
	return cdnInfo
}

func cdns(targets []string) []cdnInfo {
	stopChan = make(chan bool) // 初始化停止加载动画的通道
	// 启动加载动画的协程
	go util.ShowLoading(stopChan)

	client := req.C()
	if proxy != "" {
		client.SetProxyURL(proxy)
	}
	for _, t := range targets {
		ipList := cdnCheck(t, client)
		cdnInfoList = append(cdnInfoList, ipList)
	}
	// 任务完成，停止加载动画
	stopChan <- true
	return cdnInfoList
}

func printCDNInfos(cdnInfoList []cdnInfo) {
	res := [][]string{
		{" ", "目标", "IP数量", "检测节点", "IP列表", "检测结果"},
	}

	for _, cdnInfo := range cdnInfoList {
		for index, _ := range cdnInfo.cdnURLList { //直接找最深一层，然后将要的东西都传进去即可
			res = append(res, []string{
				strconv.Itoa(index + 1), // 从 1 开始
				cdnInfo.domain,          // 目标
				strconv.Itoa(len(strings.Split(cdnInfo.reponseList[index], ","))), //响应结果通过逗号分隔，长度就是IP数量，然后将int长度转为string
				cdnInfo.address[index],     // 检测节点，即CDN URL
				cdnInfo.reponseList[index], // IP 列表
				cdnInfo.ok[index],          // 检测结果
			})
		}
	}

	util.PrettyPrint(res)
}

func init() {
	rootCmd.AddCommand(cdnCmd)
}
