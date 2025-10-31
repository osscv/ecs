# GoECS Android App

[![Build Android APK](https://github.com/oneclickvirt/ecs/actions/workflows/build-android.yml/badge.svg)](https://github.com/oneclickvirt/ecs/actions/workflows/build-android.yml)

一个基于 Fyne 框架的跨平台测试工具的 Android 应用。

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