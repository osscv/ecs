//go:build android

package embedding

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// findNativeLibraryDir 查找应用的 native library 目录
// 这个目录是系统自动管理的，包含从 APK 中提取的 .so 文件
func findNativeLibraryDir() (string, error) {
	// 方法 1: 通过可执行文件路径推断
	execPath, err := os.Executable()
	if err == nil {
		// 可执行文件通常在 /data/app/<package>-<hash>/base.apk 或 /data/app/<package>-<hash>/oat/arm64/base.odex
		// native library 通常在 /data/app/<package>-<hash>/lib/arm64/

		// 尝试找到应用根目录
		dir := execPath
		for i := 0; i < 5; i++ { // 最多向上查找5层
			dir = filepath.Dir(dir)
			libDir := filepath.Join(dir, "lib")

			// 检查 lib 目录
			if info, err := os.Stat(libDir); err == nil && info.IsDir() {
				// 检查是否包含架构子目录
				entries, err := os.ReadDir(libDir)
				if err == nil && len(entries) > 0 {
					// 如果 lib 目录有子目录（架构名），返回 lib 目录
					for _, entry := range entries {
						if entry.IsDir() {
							return libDir, nil
						}
					}
					// 如果 lib 目录直接包含 .so 文件，也返回
					return libDir, nil
				}
			}
		}
	}

	// 方法 2: 尝试标准的 Android native library 路径
	possibleBasePaths := []string{
		"/data/data/com.oneclickvirt.goecs/lib",
		"/data/app/com.oneclickvirt.goecs/lib",
	}

	for _, basePath := range possibleBasePaths {
		if info, err := os.Stat(basePath); err == nil && info.IsDir() {
			return basePath, nil
		}

		// 尝试带哈希的路径（Android 5.0+）
		parent := filepath.Dir(basePath)
		parentEntries, err := os.ReadDir(parent)
		if err == nil {
			for _, entry := range parentEntries {
				if entry.IsDir() && strings.HasPrefix(entry.Name(), "com.oneclickvirt.goecs") {
					libDir := filepath.Join(parent, entry.Name(), "lib")
					if info, err := os.Stat(libDir); err == nil && info.IsDir() {
						return libDir, nil
					}
				}
			}
		}
	}

	// 方法 3: 搜索 /data/app 目录
	dataAppDir := "/data/app"
	if entries, err := os.ReadDir(dataAppDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() && strings.Contains(entry.Name(), "com.oneclickvirt.goecs") {
				libDir := filepath.Join(dataAppDir, entry.Name(), "lib")
				if info, err := os.Stat(libDir); err == nil && info.IsDir() {
					return libDir, nil
				}
			}
		}
	}

	return "", fmt.Errorf("无法找到 native library 目录")
}

// getLibraryName 获取当前架构对应的库名称
func getLibraryName() string {
	switch runtime.GOARCH {
	case "arm64":
		return "libgoecs_arm64.so"
	case "amd64":
		return "libgoecs_amd64.so"
	case "arm":
		return "libgoecs_arm.so"
	case "386":
		return "libgoecs_386.so"
	default:
		return "libgoecs.so"
	}
}

// ExtractECSBinary 获取 ECS 二进制文件路径
// 在 Android 上，我们不需要"提取"，而是直接使用系统已安装的 native library
func ExtractECSBinary() (string, error) {
	// 获取 native library 目录
	libDir, err := findNativeLibraryDir()
	if err != nil {
		return "", fmt.Errorf("获取 native library 目录失败: %v", err)
	}

	// 尝试的文件名列表（按优先级）
	possibleNames := []string{
		"libgoecs.so",    // 通用名称
		getLibraryName(), // 带架构后缀的名称
	}

	// 尝试的子目录（Android ABI 名称）
	abiDirs := []string{
		"", // 直接在 lib 目录
	}

	// 根据架构添加 ABI 目录
	switch runtime.GOARCH {
	case "arm64":
		abiDirs = append(abiDirs, "arm64-v8a", "arm64")
	case "arm":
		abiDirs = append(abiDirs, "armeabi-v7a", "armeabi", "arm")
	case "amd64":
		abiDirs = append(abiDirs, "x86_64", "x86-64")
	case "386":
		abiDirs = append(abiDirs, "x86")
	}

	// 尝试所有可能的路径组合
	var checkedPaths []string
	for _, abiDir := range abiDirs {
		baseDir := libDir
		if abiDir != "" {
			baseDir = filepath.Join(libDir, abiDir)
		}

		for _, name := range possibleNames {
			ecsPath := filepath.Join(baseDir, name)
			checkedPaths = append(checkedPaths, ecsPath)

			if info, err := os.Stat(ecsPath); err == nil && !info.IsDir() {
				// 找到文件，确保有执行权限
				if err := os.Chmod(ecsPath, 0755); err != nil {
					// 在某些 Android 版本上可能无法修改权限，但这通常不是问题
				}
				return ecsPath, nil
			}
		}
	}

	// 未找到文件，返回详细错误信息
	return "", fmt.Errorf("找不到 ECS 二进制文件\n已检查的路径:\n  %s\n\n请确保:\n1. ECS 二进制文件已编译为 Android 版本\n2. 文件已放置在 jniLibs/%s/libgoecs.so\n3. APK 已重新打包",
		strings.Join(checkedPaths, "\n  "),
		abiDirs[1]) // 显示推荐的 ABI 目录
}

// CleanupECSBinary 清理函数
// 在 Android 上，native library 由系统管理，我们不需要清理
func CleanupECSBinary(path string) {
	// 不需要做任何事情
	// native library 由 Android 系统管理，应用卸载时会自动清理
}
