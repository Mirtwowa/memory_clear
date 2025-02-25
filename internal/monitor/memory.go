package monitor

import (
	"awesomeProject1/config"
	"awesomeProject1/internal/utils"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ProcessInfo 存储进程信息
type ProcessInfo struct {
	Name string
	PID  string
	Mem  float64 // 内存占用 (MB)
	Path string  // 进程路径
	User string  // 所属用户
}

// StartMemoryMonitor 启动内存监控
func StartMemoryMonitor() {
	// 每5秒检查一次内存占用
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkMemoryUsage()
		}
	}
}

// checkMemoryUsage 检查内存使用情况，并判断是否超过阈值
func checkMemoryUsage() {
	// 获取当前内存使用情况
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("内存获取失败:", err)
		return
	}

	// 输出当前内存使用情况
	fmt.Printf("当前内存使用: %.2f%% (已使用: %.2fMB / 总内存: %.2fMB)\n",
		v.UsedPercent, float64(v.Used)/1024/1024, float64(v.Total)/1024/1024)

	// 如果内存使用超过设定的阈值，则记录警告日志
	if v.UsedPercent > float64(config.GetConfig().MemoryThreshold) {
		utils.LogMemoryWarning(v.UsedPercent)

		// 清理不必要的进程
		KillUnnecessaryProcesses()
	}
}

// KillUnnecessaryProcesses  查找并杀死占用内存最多的用户进程
func KillUnnecessaryProcesses() {
	// 获取配置中的忽略进程列表
	ignoreList := append(config.GetConfig().ProcessIgnoreList, "explorer.exe", "dwm.exe", "taskmgr.exe", "winlogon.exe")

	// 使用 wmic 获取进程信息
	cmd := exec.Command("wmic", "process", "get", "Name,ProcessId,WorkingSetSize,ExecutablePath", "/FORMAT:LIST")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("获取进程列表失败: %v\n", err)
		return
	}

	// 解析 wmic 输出
	processes := parseWMICOutput(string(output))
	if len(processes) == 0 {
		fmt.Println("没有找到任何进程，无法进行清理操作")
		return
	}

	// 按内存占用从大到小排序
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Mem > processes[j].Mem
	})

	// 尝试终止占用内存最多的用户进程
	for _, process := range processes {
		if isProcessInIgnoreList(process.Name, ignoreList) || isSystemProcess(process.Path) {
			continue
		}

		fmt.Printf("正在尝试终止进程: %s (PID: %s, 内存占用: %.2fMB, 路径: %s)\n", process.Name, process.PID, process.Mem, process.Path)
		killCmd := exec.Command("taskkill", "/PID", process.PID, "/F")
		output, err := killCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("无法终止进程 %s (PID: %s): %v\n输出: %s\n", process.Name, process.PID, err, string(output))
		} else {
			fmt.Printf("进程 %s (PID: %s) 已终止\n", process.Name, process.PID)
		}
	}
}

// parseWMICOutput 解析 wmic 命令输出
func parseWMICOutput(output string) []ProcessInfo {
	lines := strings.Split(output, "\n")
	var processes []ProcessInfo
	var process ProcessInfo

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			// 当前进程信息解析完成，加入结果集
			if process.Name != "" && process.PID != "" {
				processes = append(processes, process)
			}
			// 重置临时变量
			process = ProcessInfo{}
			continue
		}

		// 解析键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 根据键名填充 ProcessInfo
		switch key {
		case "Name":
			process.Name = value
		case "ProcessId":
			process.PID = value
		case "WorkingSetSize":
			memUsage, err := strconv.ParseFloat(value, 64)
			if err == nil {
				process.Mem = memUsage / 1024 / 1024 // 转换为 MB
			}
		case "ExecutablePath":
			process.Path = value
		}
	}

	// 最后一条记录
	if process.Name != "" && process.PID != "" {
		processes = append(processes, process)
	}

	return processes
}

// isProcessInIgnoreList 判断进程是否在忽略列表中
func isProcessInIgnoreList(processName string, ignoreList []string) bool {
	for _, ignoreProcess := range ignoreList {
		if strings.EqualFold(processName, ignoreProcess) {
			return true
		}
	}
	return false
}

// isSystemProcess 判断进程是否为系统进程
func isSystemProcess(path string) bool {
	if path == "" {
		return true
	}
	systemDirs := []string{
		"C:\\Windows\\System32",
		"C:\\Windows",
		"C:\\Program Files",
		"C:\\Program Files (x86)",
	}
	for _, dir := range systemDirs {
		if strings.HasPrefix(strings.ToLower(path), strings.ToLower(dir)) {
			return true
		}
	}
	return false
}
