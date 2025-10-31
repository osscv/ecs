package main

import (
	"fmt"
)

// countSelectedTests 计算已选择的测试数量
func (ui *TestUI) countSelectedTests() int {
	count := 0
	if ui.basicCheck.Checked {
		count++
	}
	if ui.cpuCheck.Checked {
		count++
	}
	if ui.memoryCheck.Checked {
		count++
	}
	if ui.diskCheck.Checked {
		count++
	}
	if ui.commCheck.Checked {
		count++
	}
	if ui.unlockCheck.Checked {
		count++
	}
	if ui.securityCheck.Checked {
		count++
	}
	if ui.emailCheck.Checked {
		count++
	}
	if ui.backtraceCheck.Checked {
		count++
	}
	if ui.nt3Check.Checked {
		count++
	}
	if ui.speedCheck.Checked {
		count++
	}
	if ui.pingCheck.Checked {
		count++
	}
	return count
}

// hasSelectedTests 检查是否有选中的测试项
func (ui *TestUI) hasSelectedTests() bool {
	return ui.basicCheck.Checked ||
		ui.cpuCheck.Checked ||
		ui.memoryCheck.Checked ||
		ui.diskCheck.Checked ||
		ui.commCheck.Checked ||
		ui.unlockCheck.Checked ||
		ui.securityCheck.Checked ||
		ui.emailCheck.Checked ||
		ui.backtraceCheck.Checked ||
		ui.nt3Check.Checked ||
		ui.speedCheck.Checked ||
		ui.pingCheck.Checked
}

// updateProgress 更新进度条和状态标签
func (ui *TestUI) updateProgress(current, total int, testName string) {
	progress := float64(current) / float64(total)
	ui.progressBar.SetValue(progress)
	ui.statusLabel.SetText(fmt.Sprintf("[%d/%d] %s", current, total, testName))
}

// isCancelled 检查测试是否被取消
func (ui *TestUI) isCancelled() bool {
	select {
	case <-ui.cancelCtx.Done():
		return true
	default:
		return false
	}
}

// resetUIState 重置UI状态
func (ui *TestUI) resetUIState() {
	ui.mu.Lock()
	ui.isRunning = false
	ui.mu.Unlock()

	ui.startButton.Enable()
	ui.stopButton.Disable()
	ui.clearButton.Enable()
	ui.progressBar.Hide()
	ui.progressBar.SetValue(0)
}
