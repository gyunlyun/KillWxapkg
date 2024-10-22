package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Ackites/KillWxapkg/cmd"
	hook2 "github.com/Ackites/KillWxapkg/internal/hook"
	"github.com/Ackites/KillWxapkg/internal/pack"
	"github.com/Ackites/KillWxapkg/scan"
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
	scanMode   bool
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
	fmt.Println("使用方法: program -id=<AppID> -in=<输入文件1,输入文件2> 或 -in=<输入目录> -out=<输出目录> [-ext=<文件后缀>] [-depth=<搜索深度>] [-restore] [-pretty] [-noClean] [-hook] [-save] [-repack=<输入目录>] [-watch] [-sensitive] [-scan]")
	flag.PrintDefaults()
	fmt.Println()
}

func init() {
	flag.StringVar(&appID, "id", "", "微信小程序的AppID")
	flag.StringVar(&input, "in", "", "输入文件路径（多个文件用逗号分隔）或输入目录路径")
	flag.StringVar(&outputDir, "out", "", "输出目录路径（如果未指定，则默认保存到输入目录下以AppID命名的文件夹）")
	flag.StringVar(&fileExt, "ext", ".wxapkg", "处理的文件后缀")
	flag.IntVar(&depth, "depth", 2, "搜索目录的深度")
	flag.BoolVar(&restoreDir, "restore", false, "是否还原工程目录结构")
	flag.BoolVar(&pretty, "pretty", false, "是否美化输出")
	flag.BoolVar(&noClean, "noClean", false, "是否清理中间文件")
	flag.BoolVar(&hook, "hook", false, "是否开启动态调试")
	flag.BoolVar(&save, "save", false, "是否保存解密后的文件")
	flag.StringVar(&repack, "repack", "", "重新打包wxapkg文件")
	flag.BoolVar(&watch, "watch", false, "是否监听将要打包的文件夹，并自动打包")
	flag.BoolVar(&sensitive, "sensitive", false, "是否获取敏感数据")
	flag.BoolVar(&scanMode, "scan", false, "扫描系统WeChat小程序目录")
}

// 新增函数：处理单个路径
func processPath(path string, currentAppID string) error {
	fmt.Printf("\n处理路径: %s\n", path)

	// 如果没有指定AppID，尝试从路径获取
	pathAppID := currentAppID
	if pathAppID == "" {
		pathAppID = cmd.GetAppID(path)
		if pathAppID == "" {
			return fmt.Errorf("无法获取路径的AppID: %s", path)
		}
	}

	// 执行处理命令
	cmd.Execute(pathAppID, path, outputDir, fileExt, depth, restoreDir, pretty, noClean, save, sensitive)
	return nil
}

func main() {
	flag.Parse()

	banner := []string{"Wxapkg", "Decompiler", "Tool", "v2.4.1"}
	colorMap := uniqueColorAssignment(banner)

	for _, word := range banner {
		fmt.Printf("%s%s\033[0m ", colorMap[word], word)
	}
	fmt.Println()

	if scanMode && (hook || watch || fileExt != ".wxapkg") {
		fmt.Println("Error: 参数冲突 - `-scan` 模式不能与 `-ext`、`-hook` 或 `-watch` 同时使用")
		printUsage()
		os.Exit(1)
	}

	// 处理scan模式
	if scanMode {
		pretty = false // 设置pretty为true
		scanner := scan.NewScanner()

		if err := scanner.ScanApplets(); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		scanner.PrintAppletPaths()

		// 循环处理每个找到的路径
		paths := scanner.GetAppletPaths()
		for i, path := range paths {
			fmt.Printf("\n[%d/%d] 开始处理路径...\n", i+1, len(paths))
			if err := processPath(path, appID); err != nil {
				fmt.Printf("处理路径时出错: %v\n", err)
				// 继续处理下一个路径，而不是退出
				continue
			}
		}
		return
	}

	// 常规模式处理
	if hook {
		hook2.Hook()
		return
	}

	if repack != "" {
		pack.Repack(repack, watch, outputDir)
		return
	}

	// 处理单个输入路径
	if input == "" {
		printUsage()
		return
	}

	if err := processPath(input, appID); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
