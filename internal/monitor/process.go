package monitor

import (
	"awesomeProject1/internal/utils"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"log"
	"time"
)

// Process 结构体定义了一个进程的基本信息
type Process struct {
	PID        int32   // 进程ID
	Name       string  // 进程名称
	CPUPercent float64 // CPU占用百分比
	MemPercent float64 // 内存占用百分比
}

// NewProcess 用于创建一个新的 Process 实例
func NewProcess(pid int32) (*Process, error) {
	// 获取进程对象
	proc, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("无法获取进程信息: %w", err)
	}

	// 获取进程的名称
	name, err := proc.Name()
	if err != nil {
		return nil, fmt.Errorf("无法获取进程名称: %w", err)
	}

	// 获取进程的 CPU 占用百分比
	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		return nil, fmt.Errorf("无法获取进程 CPU 占用: %w", err)
	}

	// 获取进程的内存占用百分比
	memPercent, err := proc.MemoryPercent()
	if err != nil {
		return nil, fmt.Errorf("无法获取进程内存占用: %w", err)
	}

	// 创建并返回一个新的 Process 实例
	return &Process{
		PID:        pid,
		Name:       name,
		CPUPercent: cpuPercent,
		MemPercent: float64(memPercent),
	}, nil
}

// Update 更新进程的 CPU 和内存占用信息
func (p *Process) Update() error {
	proc, err := process.NewProcess(p.PID)
	if err != nil {
		return fmt.Errorf("无法获取进程信息: %w", err)
	}

	// 更新 CPU 占用百分比
	cpuPercent, err := proc.CPUPercent()
	if err != nil {
		return fmt.Errorf("无法获取进程 CPU 占用: %w", err)
	}
	p.CPUPercent = cpuPercent

	// 更新内存占用百分比
	memPercent, err := proc.MemoryPercent()
	if err != nil {
		return fmt.Errorf("无法获取进程内存占用: %w", err)
	}
	p.MemPercent = float64(memPercent)

	return nil
}

// String 返回进程的简要信息
func (p *Process) String() string {
	return fmt.Sprintf("进程: %s, PID: %d, CPU占用: %.2f%%, 内存占用: %.2f%%",
		p.Name, p.PID, p.CPUPercent, p.MemPercent)
}

// StartProcessMonitor 启动进程监控
func StartProcessMonitor() {
	// 每5秒检查一次进程状态
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkProcesses()
		}
	}
}

// checkProcesses 检查系统中的所有进程，判断是否有卡死的进程
func checkProcesses() {
	// 获取系统中所有进程
	procs, err := process.Processes()
	if err != nil {
		log.Printf("获取进程列表失败: %v", err)
		return
	}

	// 遍历所有进程，检查它们的状态
	for _, proc := range procs {
		// 获取进程名称
		name, err := proc.Name()
		if err != nil {
			log.Printf("获取进程名称失败: %v", err)
			continue
		}

		// 获取进程的CPU和内存占用情况
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			log.Printf("获取进程CPU占用失败: %v", err)
			continue
		}
		memPercent, err := proc.MemoryPercent()
		if err != nil {
			log.Printf("获取进程内存占用失败: %v", err)
			continue
		}

		// 输出进程的信息
		log.Printf("进程: %s, CPU占用: %.2f%%, 内存占用: %.2f%%", name, cpuPercent, memPercent)

		// 检查进程是否处于卡死状态（根据CPU或内存占用情况）
		if isDeadlocked(cpuPercent, float64(memPercent), name) {
			utils.LogProcessWarning(name, cpuPercent, float64(memPercent))
			// 可以在此处添加自动杀死卡死进程的功能，或进一步处理
			terminateProcess(proc)
		}
	}
}

// isDeadlocked 判断进程是否处于卡死状态
// 这里简单假设：如果进程CPU占用为0且内存占用非常高，则认为进程可能卡死
func isDeadlocked(cpuPercent float64, memPercent float64, procName string) bool {
	if cpuPercent < 1.0 && memPercent > 80.0 { // 这里可以根据实际情况调整阈值
		fmt.Printf("进程 %s 可能卡死 (CPU: %.2f%%, 内存: %.2f%%)", procName, cpuPercent, memPercent)
		return true
	}
	return false
}

// terminateProcess 尝试终止卡死进程
func terminateProcess(proc *process.Process) {
	// 尝试杀死进程
	err := proc.Kill()
	if err != nil {
		fmt.Printf("无法终止进程 %d: %v", proc.Pid, err)
	} else {
		fmt.Printf("进程 %d 已终止", proc.Pid)
	}
}
