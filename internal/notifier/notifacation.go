package notifier

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// 定义缓存内存阈值 500MB
const CacheThreshold = 500 // MB

// checkCacheUsage 检查进程缓存占用情况
func CheckCacheUsage() {
	// 获取所有进程
	processes, err := getProcesses()
	if err != nil {
		log.Fatalf("无法获取进程列表: %v", err)
		return
	}

	// 遍历进程，检查缓存占用
	for _, process := range processes {
		if process.MemoryUsage > CacheThreshold {
			// 如果进程占用的内存大于阈值，发出桌面通知并让用户选择是否清除
			//notifyUser(process)
			killProcess(process)
		}
	}
}

// getProcesses 获取系统中的所有进程
func getProcesses() ([]Process, error) {
	// 在 Windows 上使用 tasklist 获取进程信息
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var processes []Process
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line != "" {
			fields := strings.Split(line, ",")
			if len(fields) >= 5 {
				processName := strings.Trim(fields[0], "\"")
				memUsageStr := strings.Trim(fields[4], "\"")
				memUsage, err := parseMemoryUsage(memUsageStr)
				if err != nil {
					continue
				}
				processes = append(processes, Process{
					Name:        processName,
					MemoryUsage: memUsage,
				})
			}
		}
	}
	return processes, nil
}

// parseMemoryUsage 解析进程的内存占用情况
func parseMemoryUsage(memUsageStr string) (float64, error) {
	memUsageStr = strings.ReplaceAll(memUsageStr, ",", "") // 去掉千位分隔符
	memUsageStr = strings.TrimSpace(memUsageStr)

	memUsage, err := strconv.ParseFloat(memUsageStr, 64)
	if err != nil {
		return 0, err
	}
	return memUsage / 1024, nil // 转换为 MB
}

// notifyUser 显示桌面通知并提示用户是否清除进程
/*func notifyUser(process Process) {
	// 启动 Fyne 应用
	myApp := app.New()

	// 创建窗口
	win := myApp.NewWindow("进程警告")

	// 创建通知内容
	content := fmt.Sprintf("进程 %s 占用了 %.2f MB 内存, 是否要清除？", process.Name, process.MemoryUsage)
	println(content)
	// 创建按钮并绑定点击事件
	agreeButton := widget.NewButton("同意", func() {
		// 用户点击同意时，清除进程
		killProcess(process)
		win.Close() // 关闭窗口
	})

	// 创建按钮并绑定点击事件
	disagreeButton := widget.NewButton("取消", func() {
		// 用户点击取消时，什么都不做
		win.Close() // 关闭窗口
	})

	// 创建按钮布局
	buttons := container.NewHBox(agreeButton, disagreeButton)

	// 创建通知对话框
	dialog.ShowCustom("内存占用警告", "确定", buttons, win)

	// 显示窗口并等待用户交互
	win.ShowAndRun()
}*/

// killProcess 结束指定进程
func killProcess(process Process) {
	cmd := exec.Command("taskkill", "/IM", process.Name, "/F")
	err := cmd.Run()
	if err != nil {
		log.Printf("无法结束进程 %s: %v", process.Name, err)
	} else {
		log.Printf("进程 %s 已被清除", process.Name)
	}
}

// Process 结构体用于存储进程信息
type Process struct {
	Name        string  // 进程名称
	MemoryUsage float64 // 内存占用（MB）
}
