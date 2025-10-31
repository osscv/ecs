package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"fyne.io/fyne/v2/dialog"
	"github.com/oneclickvirt/ecs/utils"
)

// runTests 运行所有已选择的测试
func (ui *TestUI) runTests() {
	defer func() {
		if r := recover(); r != nil {
			ui.appendResult(fmt.Sprintf("\n错误: %v\n", r))
		}
		ui.resetUIState()
	}()

	totalTests := ui.countSelectedTests()
	currentTest := 0

	language := "zh"
	if ui.languageSelect.Selected == "English" {
		language = "en"
	}

	startTime := time.Now()

	var tempOutput string

	// 打印头部并立即显示
	output := printAndCaptureGUI(func() {
		utils.PrintHead(language, width, ecsVersion)
	}, tempOutput, "")
	ui.appendResult(output)

	// 网络检测
	preCheck := utils.CheckPublicAccess(3 * time.Second)

	// 存储各项测试结果和完成状态
	var wg1, wg2, wg3 sync.WaitGroup
	var basicInfo, securityInfo string
	var mediaInfo, emailInfo, ptInfo string

	// === 第一阶段：顺序执行同步测试并实时显示 ===

	// 基础信息测试
	if ui.basicCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "基础信息测试")
		output := ui.runBasicTestCapture(language, preCheck, &basicInfo, &securityInfo, "", tempOutput)
		ui.appendResult(output)
	}

	// CPU 测试
	if ui.cpuCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "CPU 性能测试")
		output := ui.runCPUTestCapture(language, "", tempOutput)
		ui.appendResult(output)
	}

	// 内存测试
	if ui.memoryCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "内存性能测试")
		output := ui.runMemoryTestCapture(language, "", tempOutput)
		ui.appendResult(output)
	}

	// 磁盘测试
	if ui.diskCheck.Checked && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "磁盘性能测试")
		output := ui.runDiskTestCapture(language, "", tempOutput)
		ui.appendResult(output)
	}

	// === 第二阶段：启动异步测试（并发执行，但按顺序等待和显示） ===

	if ui.unlockCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			mediaInfo = utils.MediaTest(language)
		}()
	}

	if ui.emailCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			emailInfo = utils.EmailCheck()
		}()
	}

	if ui.pingCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			ptInfo = utils.PingTest()
		}()
	}

	// === 第三阶段：按顺序等待并实时显示测试结果 ===

	// 御三家流媒体测试（同步，立即显示）
	if ui.commCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "御三家流媒体测试")
		output := ui.runCommMediaTestCapture(language, "", tempOutput)
		ui.appendResult(output)
	}

	// 跨国流媒体解锁测试（等待异步完成后立即显示）
	if ui.unlockCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "跨国流媒体解锁测试")
		wg1.Wait() // 等待完成
		output := ui.runUnlockTestCapture(language, &mediaInfo, "", tempOutput)
		ui.appendResult(output) // 立即显示
	}

	// IP质量检测（立即显示）
	if ui.securityCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "IP质量检测")
		output := ui.runSecurityTestCapture(language, &securityInfo, "", tempOutput)
		ui.appendResult(output)
	}

	// 邮件端口检测（等待异步完成后立即显示）
	if ui.emailCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "邮件端口检测")
		wg2.Wait() // 等待完成
		output := ui.runEmailTestCapture(language, &emailInfo, "", tempOutput)
		ui.appendResult(output) // 立即显示
	}

	// 上游及回程线路检测（立即显示）
	if ui.backtraceCheck.Checked && preCheck.Connected && runtime.GOOS != "windows" && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "上游及回程线路检测")
		output := ui.runBacktraceTestCapture(language, "", tempOutput)
		ui.appendResult(output)
	}

	// 三网回程路由检测（立即显示）
	if ui.nt3Check.Checked && preCheck.Connected && runtime.GOOS != "windows" && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "三网回程路由检测")
		output := ui.runNT3TestCapture(language, "", tempOutput)
		ui.appendResult(output)
	}

	// 三网PING值检测（等待异步完成后立即显示）
	if ui.pingCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "三网PING值检测")
		wg3.Wait() // 等待完成
		output := ui.runPingTestCapture(language, &ptInfo, "", tempOutput)
		ui.appendResult(output) // 立即显示
	}

	// 网络测速（立即显示）
	if ui.speedCheck.Checked && preCheck.Connected && !ui.isCancelled() {
		currentTest++
		ui.updateProgress(currentTest, totalTests, "网络测速")
		output := ui.runSpeedTestCapture(language, "", tempOutput)
		ui.appendResult(output)
	}

	if !ui.isCancelled() {
		// 显示结束时间
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		minutes := int(duration.Minutes())
		seconds := int(duration.Seconds()) % 60
		currentTimeStr := endTime.Format("Mon Jan 2 15:04:05 MST 2006")

		timeOutput := printAndCaptureGUI(func() {
			utils.PrintCenteredTitle("", width)
			if language == "zh" {
				fmt.Printf("花费          : %d 分 %d 秒\n", minutes, seconds)
				fmt.Printf("时间          : %s\n", currentTimeStr)
			} else {
				fmt.Printf("Cost    Time          : %d min %d sec\n", minutes, seconds)
				fmt.Printf("Current Time          : %s\n", currentTimeStr)
			}
			utils.PrintCenteredTitle("", width)
		}, tempOutput, "")

		ui.appendResult(timeOutput)
		ui.statusLabel.SetText("测试完成")
		dialog.ShowInformation("完成", "所有测试已完成！", ui.window)
	}
}
