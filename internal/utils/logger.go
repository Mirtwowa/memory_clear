package utils

import (
	"fmt"
	"log"
	"os"
)

// InitLogger 初始化日志记录
func InitLogger() {
	// 设置日志输出到文件
	file, err := os.OpenFile("memory_monitor.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// LogMemoryWarning 记录内存警告日志
func LogMemoryWarning(memUsage float64) {
	log.Printf("警告: 当前内存使用已达到 %.2f%%", memUsage)
}

// LogProcessWarning 记录进程的警告日志
func LogProcessWarning(name string, cpuPercent float64, memPercent float64) {
	// 设置警告的阈值
	const cpuThreshold = 80.0
	const memThreshold = 90.0

	// 判断进程是否超过了 CPU 或内存的阈值
	if cpuPercent > cpuThreshold || memPercent > memThreshold {
		// 打开日志文件（如果没有则创建）
		logFile, err := os.OpenFile("process_warnings.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("无法打开日志文件: %v", err)
		}
		defer logFile.Close()

		// 创建日志记录器
		logger := log.New(logFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)

		// 格式化警告信息
		message := fmt.Sprintf("进程 '%s' 超过阈值: CPU占用 %.2f%%, 内存占用 %.2f%%", name, cpuPercent, memPercent)

		// 记录警告日志
		logger.Println(message)

		// 控制台输出
		fmt.Println(message)
	}
}
