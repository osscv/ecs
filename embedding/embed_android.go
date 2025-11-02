//go:build android

package embedding

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// isValidAppLibDir 检查是否是有效的应用 lib 目录（排除系统目录）
func isValidAppLibDir(path string) bool {
	// 排除系统目录
	systemPaths := []string{
		"/system/",
		"/vendor/",
		"/apex/",
	}

	for _, sysPath := range systemPaths {
		if strings.HasPrefix(path, sysPath) {
			return false
		}
	}

	// 必须在 /data/ 下
	if !strings.HasPrefix(path, "/data/") {
		return false
	}

	return true
}

// findNativeLibraryDir 查找应用的 native library 目录
// 这个目录是系统自动管理的，包含从 APK 中提取的 .so 文件
func findNativeLibraryDir() (string, error) {
	var allAttempts []string

	// 方法 1: 通过 /proc/self/maps 查找已加载的应用共享库路径（最可靠）
	if mapsData, err := os.ReadFile("/proc/self/maps"); err == nil {
		lines := strings.Split(string(mapsData), "\n")
		for _, line := range lines {
			// 查找包含 .so 且在 /data/ 下的行
			if strings.Contains(line, ".so") && strings.Contains(line, "/data/app/") {
				parts := strings.Fields(line)
				if len(parts) >= 6 {
					soPath := parts[5]
					allAttempts = append(allAttempts, fmt.Sprintf("从 maps 找到 .so: %s", soPath))

					// 获取库目录
					libDir := filepath.Dir(soPath)
					// 向上查找到 lib 目录
					for i := 0; i < 5; i++ {
						if filepath.Base(libDir) == "lib" && isValidAppLibDir(libDir) {
							allAttempts = append(allAttempts, fmt.Sprintf("✓ 从 maps 确定 lib 目录: %s", libDir))
							return libDir, nil
						}
						parent := filepath.Dir(libDir)
						if parent == libDir || parent == "/" {
							break
						}
						libDir = parent
					}
				}
			}
		}
	}

	// 方法 2: 通过可执行文件路径推断（仅限 /data/ 路径）
	execPath, err := os.Executable()
	if err == nil && strings.HasPrefix(execPath, "/data/") {
		allAttempts = append(allAttempts, fmt.Sprintf("可执行文件路径: %s", execPath))

		// 向上查找到包含 lib 目录的层级
		dir := execPath
		for i := 0; i < 10; i++ {
			dir = filepath.Dir(dir)
			if dir == "/" || dir == "." {
				break
			}

			libDir := filepath.Join(dir, "lib")
			allAttempts = append(allAttempts, fmt.Sprintf("检查: %s", libDir))

			if info, err := os.Stat(libDir); err == nil && info.IsDir() && isValidAppLibDir(libDir) {
				// 确保这是应用的 lib 目录（包含架构子目录）
				entries, err := os.ReadDir(libDir)
				if err == nil && len(entries) > 0 {
					allAttempts = append(allAttempts, fmt.Sprintf("✓ 从可执行路径找到 lib 目录: %s", libDir))
					return libDir, nil
				}
			}
		}
	} else {
		allAttempts = append(allAttempts, fmt.Sprintf("获取可执行路径失败或不在 /data/: %v", err))
	}

	// 方法 3: 搜索 /data/app 目录（直接扫描）
	dataAppDir := "/data/app"
	if entries, err := os.ReadDir(dataAppDir); err == nil {
		allAttempts = append(allAttempts, fmt.Sprintf("扫描 %s 目录...", dataAppDir))
		for _, entry := range entries {
			if entry.IsDir() && strings.Contains(entry.Name(), "com.oneclickvirt.goecs") {
				libDir := filepath.Join(dataAppDir, entry.Name(), "lib")
				allAttempts = append(allAttempts, fmt.Sprintf("检查: %s", libDir))
				if info, err := os.Stat(libDir); err == nil && info.IsDir() {
					allAttempts = append(allAttempts, fmt.Sprintf("✓ 从 /data/app 扫描找到 lib 目录: %s", libDir))
					return libDir, nil
				}
			}
		}
	} else {
		allAttempts = append(allAttempts, fmt.Sprintf("无法读取 %s: %v", dataAppDir, err))
	}

	// 方法 4: 尝试标准路径
	possibleBasePaths := []string{
		"/data/data/com.oneclickvirt.goecs/lib",
		"/data/app/com.oneclickvirt.goecs/lib",
	}

	for _, basePath := range possibleBasePaths {
		allAttempts = append(allAttempts, fmt.Sprintf("检查标准路径: %s", basePath))
		if info, err := os.Stat(basePath); err == nil && info.IsDir() {
			allAttempts = append(allAttempts, fmt.Sprintf("✓ 找到标准路径: %s", basePath))
			return basePath, nil
		}

		// 尝试带哈希的路径
		parent := filepath.Dir(basePath)
		if parentEntries, err := os.ReadDir(parent); err == nil {
			for _, entry := range parentEntries {
				if entry.IsDir() && strings.HasPrefix(entry.Name(), "com.oneclickvirt.goecs") {
					libDir := filepath.Join(parent, entry.Name(), "lib")
					allAttempts = append(allAttempts, fmt.Sprintf("检查带哈希的路径: %s", libDir))
					if info, err := os.Stat(libDir); err == nil && info.IsDir() {
						allAttempts = append(allAttempts, fmt.Sprintf("✓ 找到带哈希的路径: %s", libDir))
						return libDir, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("无法找到 native library 目录\n查找过程:\n  %s", strings.Join(allAttempts, "\n  "))
}

// ExtractECSBinary 获取 ECS 二进制文件路径
// 在 Android 上，我们不需要"提取"，而是直接使用系统已安装的 native library
func ExtractECSBinary() (string, error) {
	// 获取 native library 目录
	libDir, err := findNativeLibraryDir()
	debugInfo := fmt.Sprintf("架构: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	if err != nil {
		debugInfo += fmt.Sprintf("查找 lib 目录失败: %v\n", err)
	} else {
		debugInfo += fmt.Sprintf("找到 lib 目录: %s\n", libDir)

		// 列出 lib 目录内容
		if entries, err := os.ReadDir(libDir); err == nil {
			debugInfo += fmt.Sprintf("lib 目录内容 (%d 项):\n", len(entries))
			for _, entry := range entries {
				entryType := "文件"
				if entry.IsDir() {
					entryType = "目录"
				}
				debugInfo += fmt.Sprintf("  - %s (%s)\n", entry.Name(), entryType)
			}
		}
	}

	// 可能的库名称（考虑不同的命名约定）
	libraryNames := []string{
		"libgoecs.so", // 标准小写
		"libGoECS.so", // Fyne 可能使用的驼峰命名
		"libGOECS.so", // 全大写变体
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

	// 尝试所有可能的路径和名称组合
	var checkedPaths []string

	if err == nil {
		for _, abiDir := range abiDirs {
			baseDir := libDir
			if abiDir != "" {
				baseDir = filepath.Join(libDir, abiDir)
			}

			// 尝试所有可能的文件名
			for _, libraryName := range libraryNames {
				ecsPath := filepath.Join(baseDir, libraryName)
				checkedPaths = append(checkedPaths, ecsPath)

				// 添加详细的调试信息
				if info, err := os.Stat(ecsPath); err == nil && !info.IsDir() {
					// 找到文件，确保有执行权限
					if err := os.Chmod(ecsPath, 0755); err != nil {
						// 在某些 Android 版本上可能无法修改权限，但这通常不是问题
					}
					debugInfo += fmt.Sprintf("✓ 找到文件: %s\n", ecsPath)
					return ecsPath, nil
				} else if err == nil && info.IsDir() {
					debugInfo += fmt.Sprintf("  警告: %s 是目录而不是文件\n", ecsPath)
				} else {
					debugInfo += fmt.Sprintf("  未找到: %s (错误: %v)\n", ecsPath, err)
				}
			}

			// 如果这是一个目录，列出其内容
			if abiDir != "" {
				if entries, err := os.ReadDir(baseDir); err == nil && len(entries) > 0 {
					debugInfo += fmt.Sprintf("  %s 目录内容:\n", baseDir)
					for _, entry := range entries {
						entryType := "文件"
						if entry.IsDir() {
							entryType = "目录"
						}
						debugInfo += fmt.Sprintf("    - %s (%s)\n", entry.Name(), entryType)
					}
				}
			}
		}
	}

	// 如果上述方法都失败，尝试在常见位置查找
	// 注意：不再查找 /system/lib，因为那是系统库位置
	fallbackPaths := []string{
		"/data/local/tmp/libgoecs.so", // 临时目录（需要 root）
	}

	for _, path := range fallbackPaths {
		checkedPaths = append(checkedPaths, path)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}

	// 未找到文件，返回详细错误信息
	recommendedABI := "arm64-v8a"
	if runtime.GOARCH == "amd64" || runtime.GOARCH == "386" {
		recommendedABI = "x86_64"
	}

	return "", fmt.Errorf("找不到 ECS 二进制文件\n\n调试信息:\n%s\n已检查的路径:\n  %s\n\n请确保:\n1. ECS 二进制文件已编译为 Android 版本（Linux/%s）\n2. 文件已放置在 jniLibs/%s/libgoecs.so\n3. APK 已重新打包\n4. 当前架构: %s",
		debugInfo,
		strings.Join(checkedPaths, "\n  "),
		runtime.GOARCH,
		recommendedABI,
		runtime.GOARCH)
}

// CleanupECSBinary 清理函数
// 在 Android 上，native library 由系统管理，我们不需要清理
func CleanupECSBinary(path string) {
	// 不需要做任何事情
	// native library 由 Android 系统管理，应用卸载时会自动清理
}
