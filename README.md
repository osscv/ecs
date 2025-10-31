# GoECS Android App

GoECS 服务器性能测试工具的 Android 版本。

## 特性

- ✅ 图形化界面
- ✅ 7 种测试项目（基础信息、CPU、内存、磁盘、网络、流媒体、路由）
- ✅ 后台执行 + 实时进度
- ✅ 结果导出

## 快速开始

### 桌面调试
```bash
go run .
```

### 构建 APK
```bash
fyne package -os android -appID com.oneclickvirt.goecs -name GoECS
```

### 自动构建
推送到此分支会自动触发 CI 构建 APK。

## 技术栈

- Fyne v2.4.5 (UI 框架)
- Go 1.21+
- 远程依赖: github.com/oneclickvirt/ecs

## 许可证

遵循 https://github.com/oneclickvirt/ecs 项目的许可证。
