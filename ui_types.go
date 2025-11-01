package main

import (
	"context"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TestUI 测试界面结构体
type TestUI struct {
	app    fyne.App
	window fyne.Window

	// 测试选项复选框 - 完整支持所有测试项
	basicCheck     *widget.Check
	cpuCheck       *widget.Check
	memoryCheck    *widget.Check
	diskCheck      *widget.Check
	commCheck      *widget.Check // 御三家流媒体
	unlockCheck    *widget.Check // 跨国流媒体解锁
	securityCheck  *widget.Check // IP质量检测
	emailCheck     *widget.Check // 邮件端口检测
	backtraceCheck *widget.Check // 上游及回程线路
	nt3Check       *widget.Check // 三网回程路由
	speedCheck     *widget.Check // 网络测速
	pingCheck      *widget.Check // 三网PING值
	logCheck       *widget.Check // 启用日志记录

	// 预设模式选择
	presetSelect *widget.Select

	// 配置选项
	languageSelect     *widget.Select
	cpuMethodSelect    *widget.Select
	memoryMethodSelect *widget.Select
	diskMethodSelect   *widget.Select
	diskPathEntry      *widget.Entry
	threadModeSelect   *widget.Select
	nt3LocationSelect  *widget.Select
	nt3TypeSelect      *widget.Select
	diskMultiCheck     *widget.Check
	spNumEntry         *widget.Entry

	// 控制按钮
	startButton *widget.Button
	stopButton  *widget.Button

	// 结果显示 - 使用终端输出组件
	terminal    *TerminalOutput
	progressBar *widget.ProgressBar
	statusLabel *widget.Label

	// 日志相关
	logViewer *widget.Entry      // 日志查看器
	logTab    *fyne.Container    // 日志标签页内容
	mainTabs  *container.AppTabs // 主标签页容器

	// 运行状态
	isRunning bool
	cancelCtx context.Context
	cancelFn  context.CancelFunc
	mu        sync.Mutex
}
