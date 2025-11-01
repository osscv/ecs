package main

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
	ui.progressBar.Hide()
	ui.progressBar.SetValue(0)
}
