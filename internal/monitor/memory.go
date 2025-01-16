package monitor

import (
	"awesomeProject1/config"
	"awesomeProject1/internal/utils"
	"encoding/csv"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

// ProcessInfo 存储进程信息
type ProcessInfo struct {
	Name string
	PID  string
	Mem  float64 // 内存占用 (MB)
}

// 定义 Windows API 函数
var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

// getActiveProcessName 获取当前用户正在使用的程序名称
func getActiveProcessName() string {
	// 获取前台窗口句柄
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return ""
	}

	// 获取进程 ID
	var pid uint32
	_, _, _ = procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

	// 通过 tasklist 获取进程名称
	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	// 解析进程名称
	reader := csv.NewReader(strings.NewReader(string(output)))
	records, err := reader.ReadAll()
	if err != nil || len(records) == 0 || len(records[0]) < 1 {
		return ""
	}
	return strings.Trim(records[0][0], "\"")
}

// StartMemoryMonitor 启动内存监控
func StartMemoryMonitor() {
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
	v, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("内存获取失败:", err)
		return
	}

	fmt.Printf("当前内存使用: %.2f%% (已使用: %.2fMB / 总内存: %.2fMB)\n",
		v.UsedPercent, float64(v.Used)/1024/1024, float64(v.Total)/1024/1024)

	if v.UsedPercent > float64(config.GetConfig().MemoryThreshold) {
		utils.LogMemoryWarning(v.UsedPercent)
		killUnnecessaryProcesses()
	}
}

// killUnnecessaryProcesses 查找并杀死占用内存最多的进程
func killUnnecessaryProcesses() {
	ignoreList := config.GetConfig().ProcessIgnoreList
	activeProcess := getActiveProcessName()

	if activeProcess != "" {
		fmt.Printf("当前用户正在使用的程序: %s\n", activeProcess)
		ignoreList = append(ignoreList, activeProcess)
	}

	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("获取进程列表失败: %v\n", err)
		return
	}

	reader := csv.NewReader(strings.NewReader(string(output)))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("解析进程列表失败: %v\n", err)
		return
	}

	var processes []ProcessInfo
	for _, record := range records {
		if len(record) < 5 {
			continue
		}
		name := strings.Trim(record[0], "\"")
		pid := strings.Trim(record[1], "\"")
		memUsageStr := strings.Trim(record[4], "\"")

		memUsageStr = strings.ReplaceAll(memUsageStr, " K", "")
		memUsageStr = strings.ReplaceAll(memUsageStr, ",", "")
		memUsage, err := strconv.ParseFloat(memUsageStr, 64)
		if err != nil {
			fmt.Printf("解析内存使用失败: %s, 错误: %v\n", memUsageStr, err)
			continue
		}

		memUsageMB := memUsage / 1024
		processes = append(processes, ProcessInfo{Name: name, PID: pid, Mem: memUsageMB})
	}

	if len(processes) == 0 {
		fmt.Println("没有找到任何进程，无法进行清理操作")
		return
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Mem > processes[j].Mem
	})

	// 逐个终止进程，直到内存占用小于 60%
	for _, process := range processes {
		if isProcessInIgnoreList(process.Name, ignoreList) {
			continue
		}

		fmt.Printf("尝试终止进程: %s (PID: %s, 内存占用: %.2fMB)\n", process.Name, process.PID, process.Mem)
		killCmd := exec.Command("taskkill", "/PID", process.PID, "/F")
		output, err := killCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("无法终止进程 %s (PID: %s): %v\n输出: %s\n", process.Name, process.PID, err, string(output))
		} else {
			fmt.Printf("进程 %s (PID: %s) 已终止\n", process.Name, process.PID)
		}

		// 检查内存是否已降到安全阈值
		v, _ := mem.VirtualMemory()
		if v.UsedPercent <= 60 {
			fmt.Println("内存使用已降到安全范围，停止终止进程操作")
			break
		}
	}
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
