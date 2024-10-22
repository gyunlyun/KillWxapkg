package scan

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Scanner 结构体用于存储扫描相关的配置和结果
type Scanner struct {
	ConfigPath  string
	BasePath    string
	AppletPaths []string
}

// NewScanner 创建新的扫描器实例
func NewScanner() *Scanner {
	configPath := filepath.Join("C:", "Users", os.Getenv("USERNAME"), "AppData", "Roaming", "Tencent", "WeChat", "All Users", "config", "3ebffe94.ini")
	return &Scanner{
		ConfigPath: configPath,
	}
}

// ScanApplets 执行扫描过程
func (s *Scanner) ScanApplets() error {
	// 检查系统类型
	if runtime.GOOS != "windows" {
		return fmt.Errorf("scan功能仅支持Windows系统")
	}

	// 读取并解析配置文件
	if err := s.readConfig(); err != nil {
		return err
	}

	// 获取小程序路径
	if err := s.findAppletPaths(); err != nil {
		return err
	}

	return nil
}

// readConfig 读取配置文件并解析基础路径
func (s *Scanner) readConfig() error {
	content, err := ioutil.ReadFile(s.ConfigPath)
	if err != nil {
		return fmt.Errorf("无法读取配置文件: %v", err)
	}

	// 移除UTF-8 BOM和空白字符
	content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF})
	s.BasePath = strings.TrimSpace(string(content))

	log.Printf("获取微信文件路径: %s\n", s.BasePath)
	return nil
}

// findAppletPaths 查找所有wx开头的小程序目录
func (s *Scanner) findAppletPaths() error {
	appletPath := filepath.Join(s.BasePath, "WeChat Files", "Applet")

	entries, err := ioutil.ReadDir(appletPath)
	if err != nil {
		return fmt.Errorf("无法读取Applet目录: %v", err)
	}

	// 遍历查找wx开头的文件夹
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "wx") {
			fullPath := filepath.Join(appletPath, entry.Name())
			s.AppletPaths = append(s.AppletPaths, fullPath)
		}
	}

	if len(s.AppletPaths) == 0 {
		return fmt.Errorf("未找到任何wx开头的小程序目录")
	}

	return nil
}

// GetAppletPaths 返回找到的小程序路径
func (s *Scanner) GetAppletPaths() []string {
	return s.AppletPaths
}

// PrintAppletPaths 打印找到的小程序路径
func (s *Scanner) PrintAppletPaths() {
	log.Printf("找到以下小程序目录:")
	for i, path := range s.AppletPaths {
		fmt.Printf("[%d] %s\n", i+1, path)
	}
}
