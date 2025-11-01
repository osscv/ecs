//go:build !windows && !darwin && !linux

package embedding

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// 根据架构选择嵌入的二进制文件
var (
	//go:embed binaries/goecs-linux-arm64
	ecsBinaryARM64 []byte

	//go:embed binaries/goecs-linux-amd64
	ecsBinaryAMD64 []byte

	//go:embed binaries/goecs-linux-arm64
	ecsBinaryARM []byte
)

func getECSBinary() ([]byte, error) {
	switch runtime.GOARCH {
	case "arm64":
		if len(ecsBinaryARM64) == 0 {
			return nil, fmt.Errorf("ARM64 二进制文件未嵌入")
		}
		return ecsBinaryARM64, nil
	case "amd64":
		if len(ecsBinaryAMD64) == 0 {
			return nil, fmt.Errorf("AMD64 二进制文件未嵌入")
		}
		return ecsBinaryAMD64, nil
	case "arm":
		if len(ecsBinaryARM) == 0 {
			return nil, fmt.Errorf("ARM 二进制文件未嵌入")
		}
		return ecsBinaryARM, nil
	default:
		return nil, fmt.Errorf("不支持的架构: %s", runtime.GOARCH)
	}
}

func ExtractECSBinary() (string, error) {
	binary, err := getECSBinary()
	if err != nil {
		return "", err
	}

	tmpDir := os.TempDir()
	ecsPath := filepath.Join(tmpDir, "goecs")

	if err := os.WriteFile(ecsPath, binary, 0755); err != nil {
		return "", fmt.Errorf("写入二进制文件失败: %v", err)
	}

	return ecsPath, nil
}

func CleanupECSBinary(path string) {
	if path != "" {
		os.Remove(path)
	}
}
