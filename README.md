# GOECS GUI Version

[![Build All UI APP](https://github.com/oneclickvirt/ecs/actions/workflows/build.yml/badge.svg)](https://github.com/oneclickvirt/ecs/actions/workflows/build.yml)

A cross-platform testing tool based on the Fyne framework, supporting Android, macOS, and Windows.

## Supported Platforms

### Android
- Android 7.0 (API Level 24) or higher
- Android 13 (API Level 33) recommended for best experience
- Supported architectures: ARM64, x86_64

### macOS
- macOS 11.0 or higher
- Supported architectures: Apple Silicon (ARM64), Intel (AMD64)

### Windows
- Windows 10 or higher
- Supported architectures: ARM64, AMD64

## Local Build

### Prerequisites

1. Go 1.25.3
2. Android SDK
3. Android NDK 25.2.9519653
4. JDK 17+

### Environment Setup

```bash
# Set Android NDK path
export ANDROID_NDK_HOME=/path/to/android-ndk

# Install Fyne CLI
go install fyne.io/fyne/v2/cmd/fyne@latest
```

### Build Commands

```bash
# Preparation before Android build: ECS binary files need to be prepared first
# Compile Linux binaries from the ECS project and place them in the jniLibs directory
# See jniLibs/README.md for details

# Quick preparation command (assuming ecs project is in ../ecs)
cd ../ecs && \
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -checklinkname=0" -o goecs-linux-arm64 ./ && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -checklinkname=0" -o goecs-linux-amd64 ./ && \
cd ../goecs && \
cp ../ecs/goecs-linux-arm64 jniLibs/arm64-v8a/libgoecs.so && \
cp ../ecs/goecs-linux-amd64 jniLibs/x86_64/libgoecs.so && \
chmod 755 jniLibs/*/libgoecs.so

# Build desktop version (for quick testing)
./build.sh desktop

# Build Android APK (arm64 + x86_64)
./build.sh android

# Build macOS application (arm64 + amd64)
./build.sh macos

# Build Windows application (arm64 + amd64)
./build.sh windows

# Build all platforms
./build.sh all
```

Build artifacts will be output to the `.build/` directory.

### Build Artifacts Description

- **Android**: `.apk` files
  - `goecs-android-arm64-*.apk` - ARM64 version (physical device)
  - `goecs-android-x86_64-*.apk` - x86_64 version (emulator)

- **macOS**: `.tar.gz` archives (containing `.app` application)
  - `goecs-macos-arm64-*.tar.gz` - Apple Silicon version
  - `goecs-macos-amd64-*.tar.gz` - Intel version

- **Windows**: `.exe` executable files
  - `goecs-windows-arm64-*.exe` - ARM64 version
  - `goecs-windows-amd64-*.exe` - AMD64 version

## Development

```bash
# Clone repository
git clone https://github.com/oneclickvirt/ecs.git
cd ecs

# Switch to Android development branch
git checkout android-app

# Install dependencies
go mod download

# Run desktop version (for development testing)
go run -ldflags="-checklinkname=0" .
```

