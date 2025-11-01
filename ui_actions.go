package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// onPresetChanged 预设模式改变时的处理
func (ui *TestUI) onPresetChanged(preset string) {
	switch preset {
	case "1. 融合怪完全体(能测全测)":
		ui.setAllChecks(true)
	case "2. 极简版(系统+CPU+内存+磁盘+5测速节点)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.speedCheck.Checked = true
	case "3. 精简版(系统+CPU+内存+磁盘+基础解锁+5测速节点)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.unlockCheck.Checked = true
		ui.nt3Check.Checked = true
		ui.speedCheck.Checked = true
	case "4. 精简网络版(系统+CPU+内存+磁盘+回程+路由+5测速节点)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.backtraceCheck.Checked = true
		ui.nt3Check.Checked = true
		ui.speedCheck.Checked = true
	case "5. 精简解锁版(系统+CPU+内存+磁盘IO+御三家+常用流媒体+5测速节点)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.commCheck.Checked = true
		ui.unlockCheck.Checked = true
		ui.speedCheck.Checked = true
	case "6. 仅网络测试(IP质量+5测速节点)":
		ui.setAllChecks(false)
		ui.securityCheck.Checked = true
		ui.speedCheck.Checked = true
		ui.backtraceCheck.Checked = true
		ui.nt3Check.Checked = true
		ui.pingCheck.Checked = true
	case "7. 仅解锁测试(基础解锁+常用流媒体解锁)":
		ui.setAllChecks(false)
		ui.commCheck.Checked = true
		ui.unlockCheck.Checked = true
	case "8. 仅硬件测试(系统+CPU+内存+dd磁盘+fio磁盘)":
		ui.setAllChecks(false)
		ui.basicCheck.Checked = true
		ui.cpuCheck.Checked = true
		ui.memoryCheck.Checked = true
		ui.diskCheck.Checked = true
		ui.diskMethodSelect.Selected = "fio"
	case "9. IP质量测试(IP测试+15数据库+邮件端口)":
		ui.setAllChecks(false)
		ui.securityCheck.Checked = true
		ui.emailCheck.Checked = true
	default: // 自定义
		return
	}
	ui.refreshAllChecks()
}

// setAllChecks 设置所有测试项的选中状态
func (ui *TestUI) setAllChecks(checked bool) {
	ui.basicCheck.Checked = checked
	ui.cpuCheck.Checked = checked
	ui.memoryCheck.Checked = checked
	ui.diskCheck.Checked = checked
	ui.commCheck.Checked = checked
	ui.unlockCheck.Checked = checked
	ui.securityCheck.Checked = checked
	ui.emailCheck.Checked = checked
	ui.backtraceCheck.Checked = checked
	ui.nt3Check.Checked = checked
	ui.speedCheck.Checked = checked
	ui.pingCheck.Checked = checked
	ui.refreshAllChecks()
}

// refreshAllChecks 刷新所有测试项的显示
func (ui *TestUI) refreshAllChecks() {
	ui.basicCheck.Refresh()
	ui.cpuCheck.Refresh()
	ui.memoryCheck.Refresh()
	ui.diskCheck.Refresh()
	ui.commCheck.Refresh()
	ui.unlockCheck.Refresh()
	ui.securityCheck.Refresh()
	ui.emailCheck.Refresh()
	ui.backtraceCheck.Refresh()
	ui.nt3Check.Refresh()
	ui.speedCheck.Refresh()
	ui.pingCheck.Refresh()
}

// startTests 开始执行测试
func (ui *TestUI) startTests() {
	ui.mu.Lock()
	if ui.isRunning {
		ui.mu.Unlock()
		return
	}
	ui.isRunning = true
	ui.mu.Unlock()

	if !ui.hasSelectedTests() {
		dialog.ShowInformation("提示", "请至少选择一项测试！", ui.window)
		ui.mu.Lock()
		ui.isRunning = false
		ui.mu.Unlock()
		return
	}

	ui.startButton.Disable()
	ui.stopButton.Enable()
	ui.progressBar.Show()
	ui.statusLabel.SetText("测试运行中...")

	// 如果启用了日志，显示日志标签页
	if ui.logCheck != nil && ui.logCheck.Checked {
		ui.showLogTab()
	} else {
		ui.hideLogTab()
	}

	// 清空终端输出
	if ui.terminal != nil {
		ui.terminal.Clear()
	}

	ui.cancelCtx, ui.cancelFn = context.WithCancel(context.Background())
	go ui.runTestsWithExecutor()
} // stopTests 停止正在执行的测试
func (ui *TestUI) stopTests() {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	if ui.cancelFn != nil {
		ui.cancelFn()
	}
	ui.statusLabel.SetText("测试已停止")
	ui.terminal.AppendText("\n\n========== 测试被用户中断 ==========\n")
	ui.resetUIState()
}

// clearResults 清空测试结果
func (ui *TestUI) clearResults() {
	if ui.terminal != nil {
		ui.terminal.Clear()
	}
	ui.statusLabel.SetText("就绪")
	ui.progressBar.SetValue(0)
}

// exportResults 导出测试结果
func (ui *TestUI) exportResults() {
	var content string
	if ui.terminal != nil {
		content = ui.terminal.GetText()
	}

	if content == "" {
		dialog.ShowInformation("提示", "没有可导出的结果", ui.window)
		return
	}

	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, ui.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()

		_, err = writer.Write([]byte(content))
		if err != nil {
			dialog.ShowError(err, ui.window)
			return
		}

		dialog.ShowInformation("成功", "结果已导出到: "+writer.URI().Path(), ui.window)
	}, ui.window)
}

// showLogTab 显示日志标签页
func (ui *TestUI) showLogTab() {
	// 如果日志标签页还不存在，创建它
	if ui.logTab == nil {
		ui.logTab = ui.createLogTab()
		logTabItem := container.NewTabItem("日志查看", ui.logTab)
		ui.mainTabs.Append(logTabItem)
	}

	// 切换到日志标签页
	ui.mainTabs.SelectIndex(2) // 0=配置, 1=结果, 2=日志
}

// hideLogTab 隐藏日志标签页
func (ui *TestUI) hideLogTab() {
	// 如果有日志标签页，移除它
	if ui.logTab != nil && len(ui.mainTabs.Items) > 2 {
		ui.mainTabs.Remove(ui.mainTabs.Items[2])
		ui.logTab = nil
		ui.logViewer = nil
	}
}

// createLogTab 创建日志查看标签页
func (ui *TestUI) createLogTab() *fyne.Container {
	// 创建日志查看器
	ui.logViewer = widget.NewMultiLineEntry()
	ui.logViewer.Wrapping = fyne.TextWrapWord
	ui.logViewer.SetText("日志文件内容将在这里显示...")

	// 刷新按钮
	refreshButton := widget.NewButton("刷新日志", ui.refreshLog)

	// 清空按钮
	clearLogButton := widget.NewButton("清空显示", func() {
		if ui.logViewer != nil {
			ui.logViewer.SetText("")
		}
	})

	topBar := container.NewHBox(
		refreshButton,
		clearLogButton,
	)

	logScroll := container.NewScroll(ui.logViewer)

	return container.NewBorder(
		topBar,    // Top: 操作按钮
		nil,       // Bottom
		nil,       // Left
		nil,       // Right
		logScroll, // Center: 日志内容
	)
}

// refreshLog 刷新日志内容
func (ui *TestUI) refreshLog() {
	if ui.logViewer == nil {
		return
	}

	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		ui.logViewer.SetText("错误: 无法获取当前目录\n" + err.Error())
		return
	}

	// 查找所有 .log 文件
	logFiles, err := filepath.Glob(filepath.Join(currentDir, "*.log"))
	if err != nil {
		ui.logViewer.SetText("错误: 无法搜索日志文件\n" + err.Error())
		return
	}

	if len(logFiles) == 0 {
		ui.logViewer.SetText("当前目录下没有找到 .log 文件\n\n请确保已启用日志记录并运行测试。")
		return
	}

	// 找到最新的日志文件
	var latestLog string
	var latestTime time.Time

	for _, logFile := range logFiles {
		info, err := os.Stat(logFile)
		if err != nil {
			continue
		}
		if latestLog == "" || info.ModTime().After(latestTime) {
			latestLog = logFile
			latestTime = info.ModTime()
		}
	}

	if latestLog == "" {
		ui.logViewer.SetText("没有找到有效的日志文件")
		return
	}

	// 读取日志文件内容
	content, err := os.ReadFile(latestLog)
	if err != nil {
		ui.logViewer.SetText("错误: 无法读取日志文件 " + latestLog + "\n" + err.Error())
		return
	}

	// 显示日志内容
	logContent := "日志文件: " + filepath.Base(latestLog) + "\n"
	logContent += "修改时间: " + latestTime.Format("2006-01-02 15:04:05") + "\n"
	logContent += strings.Repeat("=", 60) + "\n\n"
	logContent += string(content)

	ui.logViewer.SetText(logContent)
}
