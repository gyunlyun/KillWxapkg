package cmd

import (
	"log"
	"sync"

	. "github.com/Ackites/KillWxapkg/internal/cmd"
	. "github.com/Ackites/KillWxapkg/internal/config"
	"github.com/Ackites/KillWxapkg/internal/restore"
)

func GetAppID(input string) string {
	// 调用 B 包中的 ParseWxid 函数
	log.Printf("缺少参数 id，从输入路径匹配.")
	appID := ParseWxid(input)
	return appID // 返回 appID
}

func Execute(appID, input, outputDir, fileExt string, depth int, restoreDir bool, pretty bool, noClean bool, save bool, sensitive bool) {
	// 存储配置
	configManager := NewSharedConfigManager()
	configManager.Set("appID", appID)
	configManager.Set("input", input)
	configManager.Set("outputDir", outputDir)
	configManager.Set("fileExt", fileExt)
	configManager.Set("depth", depth)
	configManager.Set("restoreDir", restoreDir)
	configManager.Set("pretty", pretty)
	configManager.Set("noClean", noClean)
	configManager.Set("save", save)
	configManager.Set("sensitive", sensitive)

	inputFiles := ParseInput(input, fileExt, depth)

	if len(inputFiles) == 0 {
		log.Println("未找到任何文件")
		return
	}

	// 确定输出目录
	outputDir = DetermineOutputDir(input, appID, outputDir)
	log.Printf("输出路径：%s", outputDir)

	var wg sync.WaitGroup
	for _, inputFile := range inputFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err := ProcessFile(file, outputDir, appID, save)
			if err != nil {
				log.Printf("处理文件 %s 时出错: %v\n", file, err)
			} else {
				log.Printf("成功处理文件: %s\n", file)
			}
		}(inputFile)
	}
	wg.Wait()

	// 还原工程目录结构
	restore.ProjectStructure(outputDir, restoreDir)
}
