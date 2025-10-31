# GoECS Android App

GoECS 服务器性能测试工具的 Android 版本。

## 快速开始

### 桌面调试
```bash
go run .
```

### 构建 APK
```bash
fyne package -os android -appID com.oneclickvirt.goecs -name GoECS
```

## 技术栈

- Fyne v2.4.5 (UI 框架)
- Go 1.21+
- 远程依赖: github.com/oneclickvirt/ecs
