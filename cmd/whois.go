package cmd

import (
	"LoveTools/util"
	"fmt"
	"regexp"
	"strings"

	"github.com/likexian/whois"
	"github.com/spf13/cobra"
)

// 定义结构化结果
type DomainInfo struct {
	Domain      string
	Registrar   string   // 注册商
	CreatedDate string   // 创建时间
	ExpiryDate  string   // 过期时间
	NameServers []string // DNS服务器
	Registrant  string   // 注册人/组织
	Status      []string // 域名状态
	UpdatedDate string   // 最后更新时间
	DNSSEC      string   // DNSSEC状态
	Emails      []string // 关联邮箱（钓鱼检测用）
}

var whoisCmd = &cobra.Command{
	Use:   "whois",
	Short: "whois 查询",
	Long:  `whois 查询，支持域名和 IP 查询`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(targets) == 0 {
			fmt.Println("缺少参数，请使用 -t 或 --target 指定目标")
			return
		}
		whoisSearch(targets)
	},
}

func whoisSearch(targets []string) {
	stopChan = make(chan bool) // 初始化停止加载动画的通道
	// 启动加载动画的协程
	go util.ShowLoading(stopChan)
	var res []string
	for _, t := range targets {
		tmp, err := whois.Whois(t)
		if err != nil {
			panic(err)
		}
		res = append(res, tmp)
	}
	// 任务完成，停止加载动画
	stopChan <- true
	// 解析WHOIS数据
	infos := parseWhois_s(res)
	// 打印结果
	printDomainInfos(infos)
}

// 解析WHOIS数据
func parseWhois(raw string) DomainInfo {

	//结构体默认值
	info := DomainInfo{
		Domain:      "nil",
		Registrar:   "nil",
		CreatedDate: "nil",
		ExpiryDate:  "nil",
		NameServers: []string{"nil"},
		Registrant:  "nil",
		Status:      []string{"nil"},
		UpdatedDate: "nil",
		DNSSEC:      "nil",
		Emails:      []string{"nil"},
	}

	// 提取域名
	info.Domain = util.ExtractField(raw, `Domain Name:\s+(.+)`)

	// 注册商（兼容不同格式）
	info.Registrar = util.FirstNonEmpty(
		util.ExtractField(raw, `Registrar:\s+(.+)`),
		util.ExtractField(raw, `Registrar Name:\s+(.+)`),
	)

	// 时间解析（自动转换时区）
	info.CreatedDate = util.ParseDate(raw, []string{
		`Creation Date:\s+(.+)`,
		`Created on:\s+(.+)`,
	})

	info.ExpiryDate = util.ParseDate(raw, []string{
		`Registry Expiry Date:\s+(.+)`,
		`Expiration Date:\s+(.+)`,
	})

	// DNS服务器（处理多行情况）
	info.NameServers = util.ExtractMultiField(raw, `Name Server:\s+(.+)`)

	// 注册人信息（去空格处理）
	info.Registrant = util.CleanSpace(
		util.ExtractField(raw, `Registrant (?:Name|Organization):\s+(.+)`),
	)

	// 域名状态（如 clientDeleteProhibited）
	info.Status = util.ExtractMultiField(raw, `Domain Status:\s+(.+)`)

	// 提取所有关联邮箱（用于钓鱼检测）
	info.Emails = extractEmails(raw)

	return info
}

func parseWhois_s(raws []string) []DomainInfo {
	var results []DomainInfo
	for _, raw := range raws {
		info := parseWhois(raw)
		results = append(results, info)
	}
	return results
}

// 提取邮箱（包括WHOIS隐私保护服务的伪装邮箱）
func extractEmails(raw string) []string {
	emailRe := regexp.MustCompile(`[\w\.-]+@[\w\.-]+\.\w+`)
	matches := emailRe.FindAllString(raw, -1)

	// 过滤常见隐私保护邮箱
	var filtered []string
	for _, email := range matches {
		if !strings.Contains(email, "whoisprotect") &&
			!strings.Contains(email, "contactprivacy") {
			filtered = append(filtered, email)
		}
	}
	return filtered
}

func printDomainInfo(info DomainInfo) {
	strs := [][]string{
		{"域名", "注册商", "创建时间", "过期时间", "DNS服务器", "注册人", "最后更新时间", "关联邮箱"},
		{info.Domain, info.Registrar,
			info.CreatedDate, info.ExpiryDate,
			strings.Join(info.NameServers, ", "),
			info.Registrant, info.UpdatedDate,
			strings.Join(info.Emails, ", ")},
	}
	util.PrettyPrint(strs)
	// fmt.Println("域名:", info.Domain)
	// fmt.Println("注册商:", info.Registrar)
	// fmt.Println("创建时间:", info.CreatedDate)
	// fmt.Println("过期时间:", info.ExpiryDate)
	// fmt.Println("DNS服务器:", strings.Join(info.NameServers, ", "))
	// fmt.Println("注册人:", info.Registrant)
	// // fmt.Println("域名状态:", strings.Join(info.Status, ", "))
	// fmt.Println("最后更新时间:", info.UpdatedDate)
	// // fmt.Println("DNSSEC状态:", info.DNSSEC)
	// fmt.Println("关联邮箱:", strings.Join(info.Emails, ", "))
}

func printDomainInfos(infos []DomainInfo) {
	for _, info := range infos {
		printDomainInfo(info)
	}
	fmt.Println("")
}

func init() {
	rootCmd.AddCommand(whoisCmd) // 添加 whois 命令
}
