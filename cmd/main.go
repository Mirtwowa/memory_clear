package main

import (
	"awesomeProject1/config"
	"awesomeProject1/internal/monitor"
	"awesomeProject1/internal/notifier"
	"awesomeProject1/internal/utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	// 初始化日志
	utils.InitLogger()

	// 初始化配置
	if conf := config.GetConfig(); conf == nil {
		println("配置加载失败")
	}
}

func main() {

	// 打印欢迎信息
	fmt.Println("Memory Monitor 启动中...")

	// 启动内存监控
	go monitor.StartMemoryMonitor()

	// 启动进程监控
	go monitor.StartProcessMonitor()

	// 启动通知模块（根据配置决定通知方式）
	go notifier.CheckCacheUsage()

	// 创建一个 channel 来监听操作系统的中断信号
	signalChannel := make(chan os.Signal, 1)

	// 捕获系统中断（Ctrl+C）或其他信号，保持程序运行直到收到终止信号
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	// 等待收到中断信号
	sigReceived := <-signalChannel
	fmt.Printf("收到信号: %v，程序正在退出...\n", sigReceived)

	// 执行清理操作，如果需要
	// 比如释放资源、关闭打开的文件/网络连接等

	// 程序退出
	fmt.Println("程序退出")
}
