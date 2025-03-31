package cmd

import (
	"LoveTools/util"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"github.com/spf13/cobra"
)

var (
	passive      bool // 是否启用被动子域名爆破
	active       bool // 是否启用主动子域名爆破
	showIP       bool // 是否显示子域名对应的ip
	urlBruteDict bool // 是否使用url字典
)

var sb = &cobra.Command{
	Use:   "subdomain",
	Short: "子域名爆破",
	Long:  `子域名爆破，支持被动和主动子域名爆破，默认启动被动子域名爆破`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(targets) == 0 {
			cmd.Help() // 显示帮助信息
		}
		if passive && active {
			fmt.Println("不能同时启用被动子域名爆破和主动子域名爆破")
			return
		}

		stopChan = make(chan bool) // 初始化停止加载动画的通道
		go util.ShowLoading(stopChan)
		if passive {
			//被动子域名爆破
			subDomainsFinder(targets)
		}
		if active {
			//主动子域名爆破
			bruteSubdomains(targets)
		}
		if !passive && !active {
			//默认启动被动子域名爆破
			subDomainsFinder(targets)
		}
		stopChan <- true // 停止加载动画
	},
}

// 由于lookup返回的结果有ip，
// 需要定义一个结构体来存储结果对应的子域名，
// 以后要用的时候方便取ip
type LookupResult struct {
	Subdomain string
	IP        []string
}

// 主动扫描
func bruteSubdomains(targets []string) []LookupResult {
	var (
		wg      sync.WaitGroup //定义一个等待组
		mutex   sync.Mutex     //定义一个互斥锁
		results []LookupResult //变量写在这里，我们的scanDomain就直接进行添加操作最方便简洁
	)

	//域名nslookup
	var scanDomain = func(domain string) {
		defer wg.Done()
		var err error
		res := LookupResult{Subdomain: domain}
		res.IP, err = net.LookupHost(domain)

		if err != nil {
			//表示子域名不存在，直接返回
			return
		}
		//否则就打印出来

		if showIP {
			fmt.Printf("\r\033[K%s: %v\n", res.Subdomain, strings.Join(res.IP, ",")) // 打印子域名和对应的ip列表
		} else {
			fmt.Printf("\r\033[K%s\n", res.Subdomain)
		}
		//加锁，因为要对results同一个蛋糕进行操作
		//所以会出现不同步的问题，可能会导致死锁
		//所以需要加锁
		mutex.Lock()
		results = append(results, res)
		mutex.Unlock()
	}

	var yamlCfg *util.Config
	yamlCfg = util.ParseConfig(yamlPath) //解析yaml
	if urlBruteDict {
		//字典为空，使用urlDict进行主动子域名爆破
		urlDict := util.LoadUrlDict(yamlCfg.UrlDict) //加载url字典
		scanner := bufio.NewScanner(urlDict)         //给到scanner读取
		for _, target := range targets {
			for scanner.Scan() {
				subdomain := scanner.Text()
				wg.Add(1)
				go scanDomain(subdomain + "." + target) //传入切片指针，虽然是值传递，但是是指针传递更清晰一点，直接在函数内部对results进行操作
			}
		}
	} else {
		//使用本地字典文件进行主动子域名爆破
		//主动扫描，使用本地字典文件
		file, _ := os.OpenFile(dict, os.O_RDONLY, 0666)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for _, target := range targets {
			for scanner.Scan() {
				subdomain := scanner.Text()
				wg.Add(1)
				go scanDomain(subdomain + "." + target) //传入切片指针，虽然是值传递，但是是指针传递更清晰一点，直接在函数内部对results进行操作
			}
		}
	}
	return results
}

// 被动扫描
// subfinder 提供的 SDK 接口主要是被动枚举子域名，不是通过传统的字典爆破来发现子域名。
// 因此，这里的 subdomain 函数只是简单地调用 subfinder 的 EnumerateSingleDomainWithCtx 函数来进行子域名爆破。
func subDomainFinder(target string) map[string]map[string]struct{} {
	subfinderOpts := &runner.Options{
		Threads:            10, // 设置线程数
		Timeout:            30, // 设置超时时间
		MaxEnumerationTime: 10, // 设置最大枚举时间

	}
	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		log.Fatalf("failed to create subfinder runner: %v", err)
	}
	output := &bytes.Buffer{} //接收结果的缓冲区
	var sourceMap map[string]map[string]struct{}
	// 枚举单个域名
	// 这里的 io.Writer 可以是任何实现了 io.Writer 接口的对象，例如文件、缓冲区等
	if sourceMap, err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), target, []io.Writer{output}); err != nil {
		log.Fatalf("failed to enumerate single domain: %v", err)
	}
	// 测试代码，打印结果
	// log.Println(output.String())

	return sourceMap

}

// 被动扫描多个域名
// subDomainsFinder 函数负责调用 subDomainFinder 函数进行多个目标子域名爆破“
func subDomainsFinder(results []string) {
	for _, target := range targets {
		//每一个域名都开启爆破动画加载
		stopChan = make(chan bool) // 初始化停止加载动画的通道
		go util.ShowLoading(stopChan)
		sourceMap := subDomainFinder(target) // 枚举子域名
		subDomainPrint(sourceMap)
		//有一个域名扫描任务完成，停止加载动画
		stopChan <- true
	}
}

// 打印被动扫描的结果
func subDomainPrint(sourceMap map[string]map[string]struct{}) {
	// 遍历 sourceMap，打印每个子域名和对应的源
	for subdomain, sources := range sourceMap {
		sourcesList := make([]string, 0, len(sources))
		for source := range sources {
			sourcesList = append(sourcesList, source)
		}
		fmt.Printf("\r\033[K%s %s (%d)\n", subdomain, sourcesList, len(sources))
	}
}

func init() {
	rootCmd.AddCommand(sb)                                                       // 添加 sb 命令
	sb.Flags().BoolVarP(&passive, "passive", "p", false, "启用被动子域名爆破")            // 添加被动子域名爆破参数
	sb.Flags().BoolVarP(&active, "active", "a", false, "启用主动子域名爆破")              // 添加主动子域名爆破参数
	sb.Flags().BoolVarP(&showIP, "showIP", "i", false, "显示子域名对应的ip")             // 添加显示子域名对应的ip参数
	sb.Flags().BoolVarP(&urlBruteDict, "udict", "u", false, "使用url字典，可在yaml中配置") // 添加url字典参数
}
