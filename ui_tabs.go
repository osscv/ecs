package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// createOptionsPanel 创建选项面板（测试项目 + 配置选项整合在一起）
func (ui *TestUI) createOptionsPanel() fyne.CanvasObject {
	// 预设模式选择
	ui.presetSelect = widget.NewSelect(
		[]string{
			"自定义",
			"1. 融合怪完全体(能测全测)",
			"2. 极简版(系统+CPU+内存+磁盘+5测速节点)",
			"3. 精简版(系统+CPU+内存+磁盘+基础解锁+5测速节点)",
			"4. 精简网络版(系统+CPU+内存+磁盘+回程+路由+5测速节点)",
			"5. 精简解锁版(系统+CPU+内存+磁盘IO+御三家+常用流媒体+5测速节点)",
			"6. 仅网络测试(IP质量+5测速节点)",
			"7. 仅解锁测试(基础解锁+常用流媒体解锁)",
			"8. 仅硬件测试(系统+CPU+内存+dd磁盘+fio磁盘)",
			"9. IP质量测试(IP测试+15数据库+邮件端口)",
		},
		ui.onPresetChanged,
	)
	ui.presetSelect.Selected = "自定义"

	presetSection := widget.NewCard("预设模式", "快速选择测试组合", ui.presetSelect)

	// === 测试项目复选框 ===
	ui.basicCheck = widget.NewCheck("基础信息测试", nil)
	ui.basicCheck.Checked = true

	ui.cpuCheck = widget.NewCheck("CPU 性能测试", nil)
	ui.cpuCheck.Checked = true

	ui.memoryCheck = widget.NewCheck("内存性能测试", nil)
	ui.memoryCheck.Checked = true

	ui.diskCheck = widget.NewCheck("磁盘性能测试", nil)
	ui.diskCheck.Checked = true

	ui.commCheck = widget.NewCheck("御三家流媒体测试", nil)
	ui.commCheck.Checked = false

	ui.unlockCheck = widget.NewCheck("跨国流媒体解锁测试", nil)
	ui.unlockCheck.Checked = false

	ui.securityCheck = widget.NewCheck("IP质量检测", nil)
	ui.securityCheck.Checked = false

	ui.emailCheck = widget.NewCheck("邮件端口检测", nil)
	ui.emailCheck.Checked = false

	ui.backtraceCheck = widget.NewCheck("上游及回程线路检测", nil)
	ui.backtraceCheck.Checked = false

	ui.nt3Check = widget.NewCheck("三网回程路由检测", nil)
	ui.nt3Check.Checked = false

	ui.speedCheck = widget.NewCheck("网络测速", nil)
	ui.speedCheck.Checked = false

	ui.pingCheck = widget.NewCheck("三网PING值检测", nil)
	ui.pingCheck.Checked = false

	ui.logCheck = widget.NewCheck("启用日志记录", nil)
	ui.logCheck.Checked = false

	// 全选/取消全选按钮
	selectAllBtn := widget.NewButton("全选", func() {
		ui.setAllChecks(true)
	})

	deselectAllBtn := widget.NewButton("取消全选", func() {
		ui.setAllChecks(false)
	})

	buttonRow := container.NewHBox(selectAllBtn, deselectAllBtn)

	// 测试项目分组 - 使用网格布局，每行2个
	basicTests := container.NewVBox(
		ui.basicCheck,
		ui.cpuCheck,
		ui.memoryCheck,
		ui.diskCheck,
	)

	networkTests := container.NewVBox(
		ui.speedCheck,
		ui.securityCheck,
		ui.emailCheck,
		ui.backtraceCheck,
	)

	advancedTests := container.NewVBox(
		ui.nt3Check,
		ui.pingCheck,
		ui.commCheck,
		ui.unlockCheck,
	)

	testsGrid := container.NewGridWithColumns(3,
		basicTests,
		networkTests,
		advancedTests,
	)

	testsSection := widget.NewCard("测试项目", "", container.NewVBox(
		buttonRow,
		testsGrid,
	))

	// === 配置选项 ===
	configSection := ui.createConfigSection()

	// 整合所有内容
	allContent := container.NewVBox(
		presetSection,
		testsSection,
		configSection,
	)

	return allContent
}

// createConfigSection 创建配置选项区域
func (ui *TestUI) createConfigSection() fyne.CanvasObject {
	// 语言选择
	ui.languageSelect = widget.NewSelect(
		[]string{"中文", "English"},
		func(value string) {},
	)
	ui.languageSelect.Selected = "中文"

	// CPU 配置
	ui.cpuMethodSelect = widget.NewSelect(
		[]string{"sysbench", "geekbench", "winsat"},
		func(value string) {},
	)
	ui.cpuMethodSelect.Selected = "sysbench"

	ui.threadModeSelect = widget.NewSelect(
		[]string{"single", "multi"},
		func(value string) {},
	)
	ui.threadModeSelect.Selected = "multi"

	// 内存配置
	ui.memoryMethodSelect = widget.NewSelect(
		[]string{"auto", "stream", "sysbench", "dd", "winsat"},
		func(value string) {},
	)
	ui.memoryMethodSelect.Selected = "auto"

	// 磁盘配置
	ui.diskMethodSelect = widget.NewSelect(
		[]string{"auto", "fio", "dd", "winsat"},
		func(value string) {},
	)
	ui.diskMethodSelect.Selected = "auto"

	ui.diskPathEntry = widget.NewEntry()
	ui.diskPathEntry.SetPlaceHolder("/tmp 或留空自动检测")

	ui.diskMultiCheck = widget.NewCheck("启用多磁盘检测", nil)
	ui.diskMultiCheck.Checked = false

	// NT3 配置
	ui.nt3LocationSelect = widget.NewSelect(
		[]string{"GZ", "SH", "BJ", "CD", "ALL"},
		func(value string) {},
	)
	ui.nt3LocationSelect.Selected = "GZ"

	ui.nt3TypeSelect = widget.NewSelect(
		[]string{"ipv4", "ipv6", "both"},
		func(value string) {},
	)
	ui.nt3TypeSelect.Selected = "ipv4"

	// 测速配置
	ui.spNumEntry = widget.NewEntry()
	ui.spNumEntry.SetText("2")
	ui.spNumEntry.SetPlaceHolder("每运营商测速节点数")

	// 使用表单布局更紧凑
	configForm := container.NewVBox(
		widget.NewLabel("通用配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("语言:"),
			ui.languageSelect,
		),
		ui.logCheck, // 日志选项
		widget.NewSeparator(),
		widget.NewLabel("CPU配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("测试方法:"),
			ui.cpuMethodSelect,
			widget.NewLabel("线程模式:"),
			ui.threadModeSelect,
		),
		widget.NewSeparator(),
		widget.NewLabel("内存配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("测试方法:"),
			ui.memoryMethodSelect,
		),
		widget.NewSeparator(),
		widget.NewLabel("磁盘配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("测试方法:"),
			ui.diskMethodSelect,
			widget.NewLabel("测试路径:"),
			ui.diskPathEntry,
		),
		ui.diskMultiCheck,
		widget.NewSeparator(),
		widget.NewLabel("三网回程配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("测试地点:"),
			ui.nt3LocationSelect,
			widget.NewLabel("测试类型:"),
			ui.nt3TypeSelect,
		),
		widget.NewSeparator(),
		widget.NewLabel("测速配置:"),
		container.NewGridWithColumns(2,
			widget.NewLabel("节点数/运营商:"),
			ui.spNumEntry,
		),
	)

	return widget.NewCard("详细配置", "调整测试参数", configForm)
}
