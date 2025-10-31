# GoECS Android App

[![Build Android APK](https://github.com/oneclickvirt/ecs/actions/workflows/build-android.yml/badge.svg)](https://github.com/oneclickvirt/ecs/actions/workflows/build-android.yml)

一个基于 Fyne 框架的跨平台服务器性能测试工具的 Android 应用。

## 功能特性

- 🔍 网络测速测试
- 🌐 端口连通性检查  
- 📊 融合怪测试
- 🎯 流媒体解锁测试
- 📱 原生 Android 界面

## 下载

在 [Releases](https://github.com/oneclickvirt/ecs/releases) 页面下载最新版本的 APK：

- **goecs-android-arm64-*.apk** - 适用于真实设备（推荐）
- **goecs-android-x86_64-*.apk** - 适用于 Android 模拟器

或在 [Actions](https://github.com/oneclickvirt/ecs/actions) 页面下载最新构建的开发版本。

## 系统要求

- Android 7.0 (API Level 24) 或更高版本
- 建议 Android 13 (API Level 33) 以获得最佳体验

## 本地构建

### 前置要求

1. Go 1.23+
2. Android SDK
3. Android NDK 25.2.9519653
4. JDK 17+

### 环境配置

```bash
# 设置 Android NDK 路径
export ANDROID_NDK_HOME=/path/to/android-ndk

# 安装 Fyne CLI
go install fyne.io/fyne/v2/cmd/fyne@latest
```

### 构建命令

```bash
# 构建桌面端（用于快速测试）
./build.sh desktop

# 构建 Android APK
./build.sh android

# 构建所有平台
./build.sh all
```

构建产物将输出到 `.build/` 目录。

## 自动构建

项目使用 GitHub Actions 自动构建，每次推送到 `main` 或 `android-app` 分支时都会触发构建：

- ✅ 自动编译 ARM64 和 x86_64 两个架构的 APK
- ✅ 自动上传构建产物
- ✅ 标签推送时自动创建 Release
- ✅ 自动提交构建好的 APK 到仓库

## 开发

```bash
# 克隆仓库
git clone https://github.com/oneclickvirt/ecs.git
cd ecs

# 切换到 Android 开发分支
git checkout android-app

# 安装依赖
go mod download

# 运行桌面版本（用于开发测试）
go run -ldflags="-checklinkname=0" .
```

## 技术栈

- [Fyne](https://fyne.io/) - 跨平台 GUI 框架
- [Go](https://go.dev/) - 编程语言
- [oneclickvirt/ecs](https://github.com/oneclickvirt/ecs) - 核心测试库

## License

MIT License

## 相关项目

- [oneclickvirt/ecs](https://github.com/oneclickvirt/ecs) - 服务器性能测试命令行工具
