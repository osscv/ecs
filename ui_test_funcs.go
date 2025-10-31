package main

import (
	"fmt"

	"github.com/oneclickvirt/CommonMediaTests/commediatests"
	"github.com/oneclickvirt/ecs/cputest"
	"github.com/oneclickvirt/ecs/disktest"
	"github.com/oneclickvirt/ecs/memorytest"
	"github.com/oneclickvirt/ecs/nexttrace"
	"github.com/oneclickvirt/ecs/speedtest"
	"github.com/oneclickvirt/ecs/upstreams"
	"github.com/oneclickvirt/ecs/utils"
)

// 各测试函数实现 - 只捕获输出并返回

func (ui *TestUI) runBasicTestCapture(language string, preCheck utils.NetCheckResult, basicInfo, securityInfo *string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("系统基础信息", width)
		} else {
			utils.PrintCenteredTitle("System-Basic-Information", width)
		}

		var nt3CheckType string = ui.nt3TypeSelect.Selected

		if preCheck.Connected && preCheck.StackType == "DualStack" {
			_, _, *basicInfo, *securityInfo, nt3CheckType = utils.BasicsAndSecurityCheck(language, nt3CheckType, ui.securityCheck.Checked)
		} else if preCheck.Connected && preCheck.StackType == "IPv4" {
			_, _, *basicInfo, *securityInfo, nt3CheckType = utils.BasicsAndSecurityCheck(language, "ipv4", ui.securityCheck.Checked)
		} else if preCheck.Connected && preCheck.StackType == "IPv6" {
			_, _, *basicInfo, *securityInfo, nt3CheckType = utils.BasicsAndSecurityCheck(language, "ipv6", ui.securityCheck.Checked)
		} else {
			_, _, *basicInfo, *securityInfo, nt3CheckType = utils.BasicsAndSecurityCheck(language, "", false)
		}

		fmt.Printf("%s", *basicInfo)
	}, tempOutput, output)
}

func (ui *TestUI) runCPUTestCapture(language string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		realTestMethod, res := cputest.CpuTest(language, ui.cpuMethodSelect.Selected, ui.threadModeSelect.Selected)
		if language == "zh" {
			utils.PrintCenteredTitle(fmt.Sprintf("CPU测试-通过%s测试", realTestMethod), width)
		} else {
			utils.PrintCenteredTitle(fmt.Sprintf("CPU-Test--%s-Method", realTestMethod), width)
		}
		fmt.Print(res)
	}, tempOutput, output)
}

func (ui *TestUI) runMemoryTestCapture(language string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		realTestMethod, res := memorytest.MemoryTest(language, ui.memoryMethodSelect.Selected)
		if language == "zh" {
			utils.PrintCenteredTitle(fmt.Sprintf("内存测试-通过%s测试", realTestMethod), width)
		} else {
			utils.PrintCenteredTitle(fmt.Sprintf("Memory-Test--%s-Method", realTestMethod), width)
		}
		fmt.Print(res)
	}, tempOutput, output)
}

func (ui *TestUI) runDiskTestCapture(language string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		diskPath := ui.diskPathEntry.Text
		diskMethod := ui.diskMethodSelect.Selected
		diskMultiCheck := ui.diskMultiCheck.Checked
		autoChange := (diskMethod == "auto")

		realTestMethod, res := disktest.DiskTest(language, diskMethod, diskPath, diskMultiCheck, autoChange)
		if language == "zh" {
			utils.PrintCenteredTitle(fmt.Sprintf("硬盘测试-通过%s测试", realTestMethod), width)
		} else {
			utils.PrintCenteredTitle(fmt.Sprintf("Disk-Test--%s-Method", realTestMethod), width)
		}
		fmt.Print(res)
	}, tempOutput, output)
}

func (ui *TestUI) runCommMediaTestCapture(language string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("御三家流媒体解锁", width)
		} else {
			utils.PrintCenteredTitle("Common-Streaming-Media-Unlock", width)
		}
		fmt.Printf("%s", commediatests.MediaTests(language))
	}, tempOutput, output)
}

func (ui *TestUI) runUnlockTestCapture(language string, mediaInfo *string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("跨国流媒体解锁", width)
		} else {
			utils.PrintCenteredTitle("Cross-Border-Streaming-Media-Unlock", width)
		}
		fmt.Printf("%s", *mediaInfo)
	}, tempOutput, output)
}

func (ui *TestUI) runSecurityTestCapture(language string, securityInfo *string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("IP质量检测", width)
		} else {
			utils.PrintCenteredTitle("IP-Quality-Check", width)
		}
		fmt.Printf("%s", *securityInfo)
	}, tempOutput, output)
}

func (ui *TestUI) runEmailTestCapture(language string, emailInfo *string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("邮件端口检测", width)
		} else {
			utils.PrintCenteredTitle("Email-Port-Check", width)
		}
		fmt.Println(*emailInfo)
	}, tempOutput, output)
}

func (ui *TestUI) runBacktraceTestCapture(language string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("上游及回程线路检测", width)
		} else {
			utils.PrintCenteredTitle("Upstream-and-Return-Route-Check", width)
		}
		upstreams.UpstreamsCheck()
	}, tempOutput, output)
}

func (ui *TestUI) runNT3TestCapture(language string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("三网回程路由检测", width)
		} else {
			utils.PrintCenteredTitle("Three-Network-Return-Route-Check", width)
		}
		nexttrace.NextTrace3Check(language, ui.nt3LocationSelect.Selected, ui.nt3TypeSelect.Selected)
	}, tempOutput, output)
}

func (ui *TestUI) runPingTestCapture(language string, ptInfo *string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("三网ICMP的PING值检测", width)
		} else {
			utils.PrintCenteredTitle("Three-Network-ICMP-Ping-Check", width)
		}
		fmt.Println(*ptInfo)
	}, tempOutput, output)
}

func (ui *TestUI) runSpeedTestCapture(language string, output, tempOutput string) string {
	return printAndCaptureGUI(func() {
		if language == "zh" {
			utils.PrintCenteredTitle("就近节点测速", width)
		} else {
			utils.PrintCenteredTitle("Speed-Test", width)
		}

		speedtest.ShowHead(language)
		speedtest.NearbySP()

		// 根据预设模式调整测速节点数
		spNum := 2
		if ui.spNumEntry.Text != "" {
			fmt.Sscanf(ui.spNumEntry.Text, "%d", &spNum)
		}

		speedtest.CustomSP("net", "global", spNum, language)
	}, tempOutput, output)
}
