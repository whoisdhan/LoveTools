# 个人武器开发

国内项目地址：[https://gitee.com/whoisDhan/LoveTools](https://gitee.com/whoisDhan/LoveTools)
国外项目地址：
这里可以开始给你第一个安全工具起名字了，\_\_\_\_\_\_\_ <- 名字留给各位师傅替自己填
注：本章节的所有代码放github仓库中，或者可以到文章最结尾复制全部源码，作者建议各位还是去github上直接下载整个项目来学习，不用自己创建了
（如果喜欢的话可以顺手点一个免费的star~后续我肯定会出一些安全工具的，各位师傅可以持续关注我的公众号！！！）


---


tips:
由于代码我在最后才放，所以这里先告诉各位师傅我存了两个全局参数在root.go中，target是通用的，proxy也是通用的，同时target目标可以用逗号分隔开作为多个目标传递，因为我使用了StringSliceVarP函数接收参数，OK交代完毕。
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140536945-597643523.png)
## 美观输出

先对后面的终端输出简单的美观装饰一下，装b必备，学了安全开发可以给不懂安全的朋友装一下多是一件美事。

- 设置表头：SetHeader(v[0])  
- 启用边框：SetBorder(true) 
- 启用行分隔线：SetRowLine(true)  
- 开启自动换行：SetAutoWrapText(true) 
- 自定义中心分隔符：SetCenterSeparator("+")  // 即每一个列之间的分隔符
- 自定义列分隔符：SetColumnSeparator("|")
- 自定义行分隔符：SetRowSeparator("-")

当我们传入数据调用函数的时候
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140536668-1223981883.png)

输出如下图所示：
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140536382-317664310.png)

函数代码如下：
```go
func PrettyPrint(v [][]string) {

    table := tablewriter.NewWriter(os.Stdout)

    table.SetHeader(v[0])         // 设置表头

    table.SetBorder(true)         // 启用边框

    table.SetRowLine(true)        // 启用行分隔线

    table.SetAutoWrapText(true)   // 开启自动换行

    table.SetCenterSeparator("+") // 自定义中心分隔符，即每一个列之间的分隔符

    table.SetColumnSeparator("|") // 自定义列分隔符

    table.SetRowSeparator("-")    // 自定义行分隔符

  

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
```


## Whois查询

下载
```go
go get -u github.com/likexian/whois
```

简单测试一个域名：
```go
t := "baidu.com"
tmp, err := whois.Whois(t)

if err != nil {

	panic(err)

}
fmt.Println(res)

```
这个是最简单的用法，然后为了接收返回的结果，我创建了一个结构体来接收
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140536017-1884551365.png)
但是whois返回的是整个查询结果的原生字符串，所以只能够自己写正则匹配来提取内容了
代码较多，就发核心代码截图来讲解（所有源码在结尾）

util文件的一些函数
- 加载动画
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140535722-291303939.png)
- whois主要功能
	其中包含了加载动画、正则提取、打印结果
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140535428-230836670.png)

- 运行结果
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140535116-741815072.png)

更多函数细节不用深究，后面我会给出所有源码，而且这些函数功能其实用ai也能写出来，不用造轮子。

## 反查ip

最终效果
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140534853-2090768432.png)
自己在[https://site.ip138.com/](https://site.ip138.com/)网站上找ip进行反查作为例子即可，很多ip都可以反查到域名
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140534556-1096605365.png)

---


这功能在[https://site.ip138.com/](https://site.ip138.com/)网站请求

- 核心函数`iprSearch`:
	负责对单个域名进行ip反查
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140534233-889292067.png)

- `iprsSearch`，ipr后多了一个s区分，用来对域名列表进行ip反查
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140533940-970723683.png)
- `ipr`最终调用功能函数
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140533636-413258490.png)

- 运行结果
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140533349-1590517053.png)


## 目录扫描

练习项目我们这里就直接用字典扫url一层目录即可
- 2xx 状态码：绿色
- 3xx 状态码：橙色
- 4xx 状态码：蓝色
- 5xx 状态码：红色
最终运行效果如下：
（我们仅仅做于学习，并没有对服务器造成影响，扫了几个请求而已哈~）
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140533097-1690254187.png)

细节：
- 目录扫描要设置禁止自动跟踪跳转，否则服务器有完整的302跳转的话，你扫出来的都是200
- 超时时间建议也带上，可控性强
- 这里其实还可以加一个延迟请求效果，比如每一个请求之间间隔多久，否则发送太快容易被服务端封禁（这里留着可以给感兴趣的师傅自己写）
- 线程设置，这里不做实现（×）
- 我们scanner读取每一行的时候，由于可能用户会给多个目标，所以我们的字典要指针要回到头部进行重新读取，seek需要第二层循环完成字典后记得将指针指向头部

![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140532790-971292073.png)


对url进行清洗
- checkHttp：检查是否是http开头，因为有的用户可能会直接给一个域名
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140532487-1466617690.png)
- 去除多余空格
- 去除多余`/`

![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140532156-1828549533.png)

对不同状态码之间的请求路径上色
- 这里判断一下状态码的范围即可，注意我第一个用了`IsSuccessState`，他就是判断200~299之间的，后面就没有了，只剩下一个`IsErrorState`，他是大于400就true，不符合我们的要求
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140531822-301957593.png)
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140531540-34696504.png)

一个目标完成扫描后，我们的字典到尾部了，所以我们需要将这个指针指向头部重新给下一个目标扫描目录
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140531229-1645029417.png)

最后给命令Run给上运行逻辑和init上添加子命令和对应的参数即可
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140530924-622777895.png)

最后就是运行`go run main.go -t xxx.com,aaa.ccc`  ,多个目标可以用逗号隔开
(仅仅做学习，不要对服务器造成影响)
这里还有几个状态码效果颜色需要找到对应的服务器返回状态码才行，这里就忽略了，能够打印颜色就表示成功了。
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140530535-1063835791.png)
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140530251-877710272.png)

## 子域名爆破

这里有一个细节要注意：
- 由于子域名爆破中主动扫描要用到dict字典，之前在目录扫描中也有一个字典，为了方便，使用同一个参数就合并起来了，放到root中作为全局参数
	**`同时为了更容易区分和防止子命令撞参数就把全局的短选项参数更改为大写`**
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140529929-1966990995.png)
- 同时也在root添加了yaml配置文件变量名
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140529636-536529240.png)

- 被动扫描用的是subfinder提供出来的sdk进行开发，下载的时候需要注意go get是否能下载到，下载不了可能你是更改成为了其他加速的地址需要更改回来官方的：
	**具体要看你能不能下载，下载失败就要更改回官方的**
	`go env -w GOPROXY=https://goproxy.io,direct` <- 这是官方的，之前可能你更改了阿里云或者其他国内加速地址
### 被动扫描

使用subfinder的核心接口sdk，有官方使用代码例子：
`https://github.com/projectdiscovery/subfinder/blob/dev/v2/examples/main.go`


我这里写了三个函数，将官方示例代码小小的拆开了几部分
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140529194-1678403983.png)
- `subDomainFinder`：扫描单个域名
	0.写一个`&runner.Options`结构体
	1.`NewRunner`创建一个runner
	2.`&bytes.Buffer{}`缓冲区
	3.使用`EnumerateSingleDomainWithCtx` 进行被动扫描，参数按照官方给的用就行，这些代码都是官方上拿的
	4.改动：我将结果作为函数返回值，因为我们是对多个目标进行扫描，这里单个结果返回即可
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140528880-832871663.png)

- 对多个目标进行扫描，就是结合单个域名扫描那里进行包装循环
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140528564-1479119243.png)

- 打印结果：没啥好说的，就看你自己需求，我这里就直接用官方给的打印方式了，唯一不同就是由于我们有加载动画，首先不能让加载动画覆盖我们的域名打印结果，所以我们在打印结果之前进行清行：`\r\033[K`，当然
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140528275-67462964.png)
### 主动扫描(字典爆破)


主动扫描意思是使用字典进行爆破，所以自己就能写，可以不用subfinder的代码，使用`lookup`就能判断域名是否有效。
我这里用了一个函数一个结构体
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140527961-838166394.png)

- 结构体`LookupResult`：
	用来存储域名和对应的ip，因为他`lookup`会返回一个lookup的ip地址列表

- `bruteSubdomains`：go协程、加锁、文件读取、yaml文件解析都用上了，这里学习了一个新的解析格式`yaml`，建议自己去看[https://www.runoob.com/w3cnote/yaml-intro.html](https://www.runoob.com/w3cnote/yaml-intro.html)菜鸟这篇文章，短小精悍，一看就懂。
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140527668-588345330.png)
	- 单个域名解析：`scanDomain` 函数变量
		这里还涉及到多个协程操作一个`results`列表变量的问题，所以需要用到互斥锁，否则会出现不同步的问题，可能会导致死锁。
		还有一个就是ip显示还是不显示的问题，这个也提取出来作为一个可选参数，毕竟我们爆破域名的时候也想看看ip（默认是不开启）
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140527391-619781899.png)
	- 为了更加定制化，所以我将url字典提取出来了，放到了yaml文件中，方便以后url字典能够在配置文件更换
		![](md图片/README/IMG-20250331140831587.png)
	- 这里需要添加一个config.go文件存储结构体
		![](md图片/README/IMG-20250331140831644.png)
	- 接着就是写解析`yaml`函数，其实就是之前学的导出文件中的函数一样，只不过函数类别换成了yaml来调用
		![](md图片/README/IMG-20250331140831811.png)
	- 因为要用`bufio.NewScanner`，在util.go文件中添加一个函数，所以要读取url返回一个`io.Reader`类型，那么就直接封装到`util.go`文件中使用了
		![](md图片/README/IMG-20250331140831985.png)
	- 本地扫描就忽略了，因为就是一个读取文件一行一行拼接就行，下面我直接放剩下的代码截图
		![](md图片/README/IMG-20250331140832145.png)


- 参数添加如下
	所以我们可以使用的命令组合有：
	`subdomain -a true -T xxx.cxm` 默认使用url字典主动扫描
	`subdomain -a true -F dict.txt -T xxx.cxm` 使用字典路径扫描
	`subdomain -b true -T xxx.cxm` 被动扫描
	`subdomain -T xxx.cxm`  两个都是false的时候，默认就是被动扫描，函数代码中有进行判断两个false
	`subdomain -i true -T xxx.cxm`加上-i组合就更多了，自己测试即可
![](md图片/README/IMG-20250331140832382.png)
这里贴一个运行结果
命令是：`go run main.go subdomain -a true -u true -T baidu.com -i true`
![](md图片/README/IMG-20250331140832562.png)

### CDN检测

检测最终结果，仅做学习用途，无非法动作。
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140527105-2128309964.png)

- yaml文件添加cdn检测节点
	节点检测方式：`CDNURL/?domain=baidu.com`  -> 响应返回用逗号隔开的ip列表(检测成功) 或者检测失败
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140526737-223395807.png)

- config.go文件添加属性
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140526445-1218784888.png)

- 创建`cdn.go`文件，分别写了三个函数以及一个结构体去完成这个功能
	- `cdnInfo`就是用来存储一个域名去多个cdn节点检测结果的一个结果合集
		所以除了domain都是string列表类型，看下图抓包就能看到请求的完整路径，所以知道请求路径就知道怎么写代码了。
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140526130-1657657550.png)
![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140525754-957372229.png)
- `cdnCheck`最核心的函数：最终要的是在你`config := util.ParseConfig(yamlPath)`读取配置文件后怎么请求，其实就是很简单，用创建好的client请求，如果状态码200就请求成功，否则就请求失败，在响应结果和相应列表的append中自己修改一下即可。
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140525418-1361623492.png)
	- 在第二层循环的时候记得把`target`和yaml配置中的cdn归属地址`address`添加进`cdnInfo`里面即可

- `cdns函数`：
	就是对多个目标进行遍历给到cdnCheck检测即可
	同时记得加上加载动画过程
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140525098-31622917.png)
- 打印CDN结果
	双层遍历找到`cdnInfo`的结果列表，因为他就是最深一层，然后按照需求添加进去`res`里面给到`util.PrettyPrint`函数就行(这个函数是之前写的，一直都有用)
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140524401-992542330.png)
- Command结构体的Run属性写法：
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140524099-1533545748.png)
- 最终运行结果：
	正在扫描中
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140523762-27707923.png)
	结果
	![](https://img2023.cnblogs.com/blog/3392862/202503/3392862-20250331140523249-762065046.png)