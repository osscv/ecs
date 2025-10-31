package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// NewTestUI 创建新的测试UI实例
func NewTestUI(app fyne.App) *TestUI {
	ui := &TestUI{
		app:    app,
		window: app.NewWindow("融合怪测试 - 完整版"),
	}

	// 设置窗口大小 - 增大以适应左右分栏布局
	ui.window.Resize(fyne.NewSize(1400, 900))
	ui.window.SetPadded(true)
	ui.window.CenterOnScreen()

	ui.buildUI()
	return ui
}

// buildUI 构建用户界面
func (ui *TestUI) buildUI() {
	// 创建选项卡
	tabs := container.NewAppTabs(
		container.NewTabItem("测试项目", ui.createTestOptionsTab()),
		container.NewTabItem("配置选项", ui.createConfigTab()),
		container.NewTabItem("测试结果", ui.createResultTab()),
	)

	ui.window.SetContent(tabs)
}
