package main

import (
	"flag"
	"fmt"
	"github.com/Ackites/KillWxapkg/internal/pack"
	"math/rand"
	"time"

	"github.com/Ackites/KillWxapkg/cmd"
	hook2 "github.com/Ackites/KillWxapkg/internal/hook"
)

var colors = []string{
	"\033[32m", // 绿色
	"\033[33m", // 黄色
	"\033[34m", // 蓝色
	"\033[35m", // 紫色
	"\033[36m", // 青色
	"\033[92m", // 绿色
	"\033[93m", // 黄色
	"\033[94m", // 蓝色
	"\033[95m", // 紫色
	"\033[96m", // 青色
}

var (
	appID      string
	input      string
	outputDir  string
	fileExt    string
	restoreDir bool
	pretty     bool
	noClean    bool
	hook       bool
	save       bool
	repack     string
	watch      bool
	sensitive  bool
	depth      int
)

func uniqueColorAssignment(words []string) map[string]string {
	rand.Seed(time.Now().UnixNano())
	shuffledColors := rand.Perm(len(colors))
	colorMap := make(map[string]string)

	for i, word := range words {
		if i < len(shuffledColors) {
			colorMap[word] = colors[shuffledColors[i]]
		}
	}

	return colorMap
}

func printUsage() {
	fmt.Println("使用方法: program -id=<AppID> -in=<输入文件1,输入文件2> 或 -in=<输入目录> -out=<输出目录> [-ext=<文件后缀>] [-depth=<搜索深度>] [-restore] [-pretty] [-noClean] [-hook] [-save] [-repack=<输入目录>] [-watch] [-sensitive]")
	flag.PrintDefaults()
	fmt.Println()
}

func init() {
	flag.StringVar(&appID, "id", "", "微信小程序的AppID")
	flag.StringVar(&input, "in", "", "输入文件路径（多个文件用逗号分隔）或输入目录路径")
	flag.StringVar(&outputDir, "out", "", "输出目录路径（如果未指定，则默认保存到输入目录下以AppID命名的文件夹）")
	flag.StringVar(&fileExt, "ext", ".wxapkg", "处理的文件后缀")
	flag.IntVar(&depth, "depth", 2, "搜索目录的深度（默认为2）") // 新增 depth 参数
	flag.BoolVar(&restoreDir, "restore", false, "是否还原工程目录结构")
	flag.BoolVar(&pretty, "pretty", false, "是否美化输出")
	flag.BoolVar(&noClean, "noClean", false, "是否清理中间文件")
	flag.BoolVar(&hook, "hook", false, "是否开启动态调试")
	flag.BoolVar(&save, "save", false, "是否保存解密后的文件")
	flag.StringVar(&repack, "repack", "", "重新打包wxapkg文件")
	flag.BoolVar(&watch, "watch", false, "是否监听将要打包的文件夹，并自动打包")
	flag.BoolVar(&sensitive, "sensitive", false, "是否获取敏感数据")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	banner := []string{"Wxapkg", "Decompiler", "Tool", "v2.4.1"}

	// 分配不重复的颜色
	colorMap := uniqueColorAssignment(banner)

	// 打印每个单词及其分配的颜色
	for _, word := range banner {
		fmt.Printf("%s%s\033[0m ", colorMap[word], word)
	}
	fmt.Println() // 换行

	// 动态调试
	if hook {
		hook2.Hook()
		return
	}

	// 重新打包
	if repack != "" {
		pack.Repack(repack, watch, outputDir)
		return
	}

	if input == "" {
		printUsage()
		return
	}

	appID := cmd.GetAppID(input)

	if appID == "" {
		printUsage()
		return
	}

	// 执行命令
	cmd.Execute(appID, input, outputDir, fileExt, depth, restoreDir, pretty, noClean, save, sensitive)
}
