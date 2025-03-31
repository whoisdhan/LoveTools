package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

// ShowLoading 显示加载动画
func ShowLoading(stopChan chan bool) {
	frames := []string{"|", "|", "/", "/", "-", "-", "\\", "\\"}
	dot := []string{".", ".", "..", "..", "...", "...", "....", "...."}
	i := 0
	for {
		select {
		case <-stopChan:
			fmt.Printf("\r\033[K") // Clear the line
			return
		default:
			fmt.Printf("\r\033[K[%s] Scanning%s", frames[i], dot[i])
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(frames)
		}
	}
}

// 美观输出
func PrettyPrint(v [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(v[0])         // 设置表头
	table.SetBorder(true)         // 启用边框
	table.SetRowLine(true)        // 启用行分隔线
	table.SetAutoWrapText(true)   // 开启自动换行
	table.SetCenterSeparator("+") // 自定义中心分隔符，即每一个列之间的分隔符
	table.SetColumnSeparator("|") // 自定义列分隔符
	table.SetRowSeparator("-")    // 自定义行分隔符

	// 添加数据
	if len(v) == 0 {
		//如果没有数据，直接返回
		table.Render()
		return
	}
	for _, row := range v[1:] {
		// 这里v[1:]从第二行开始添加数据，第一行是表头，因为表头已经设置好了
		table.Append(row)
	}
	table.Render() // 将结果 prettily 打印到标准输出
	fmt.Println("")
}

// 正则提取
func ExtractField(raw string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(raw)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// 提取多个匹配项
func ExtractMultiField(raw string, pattern string) []string {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(raw, -1)
	var result []string
	for _, m := range matches {
		if len(m) > 1 {
			result = append(result, strings.TrimSpace(m[1]))
		}
	}
	return result
}

// 时间解析（兼容多种格式）
func ParseDate(raw string, patterns []string) string {
	for _, pattern := range patterns {
		if value := ExtractField(raw, pattern); value != "" {
			// 尝试解析常见时间格式
			formats := []string{
				"2006-01-02T15:04:05Z",
				"2006-01-02",
				"02-Jan-2006",
			}
			for _, layout := range formats {
				if t, err := time.Parse(layout, value); err == nil {
					return t.Format("2006-01-02 15:04:05 MST")
				}
			}
			return value // 无法解析则返回原始字符串
		}
	}
	return ""
}

// 其他辅助函数
func FirstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// 合并多个空格为一个空格
// 例如： "  a   b  c  " => "a b c"
func CleanSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// 去除空格
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// 去除 / 符号
func TrimSlash(s string) string {
	return strings.Trim(s, "/")
}

// 去除尾部多个 / 符号
func TrimSlashEnd(s string) string {
	return strings.TrimRight(s, "/")
}

// 去除前后 / 符号， 中间的 / 不去除
// 例如： "/a/b/c/" => "a/b/c"
func TrimSlashProper(s string) string {
	s = strings.TrimLeft(s, "/")   // 去除开头所有 `/`
	s = strings.TrimSuffix(s, "/") // 去除结尾的 `/`
	return s
}

// 转为file一样正常给scanner读取
func LoadUrlDict(url string) io.Reader {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("请求失败url字典")
		os.Exit(1)
	}
	defer resp.Body.Close()
	// 读取 Body 数据
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取失败")
		os.Exit(1)
	}

	// 用 bytes.Reader 返回
	return bytes.NewReader(data)
}
