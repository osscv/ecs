//go:build !windows && !darwin && !android

package embedding

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed binaries/goecs-linux-arm64
var ecsBinaryARM64 []byte

//go:embed binaries/goecs-linux-amd64
var ecsBinaryAMD64 []byte

func getECSBinary() ([]byte, error) {
	switch runtime.GOARCH {
	case "arm64":
		if len(ecsBinaryARM64) == 0 {
			return nil, fmt.Errorf("ARM64 二进制文件未嵌入 (GOOS=%s, GOARCH=%s)", runtime.GOOS, runtime.GOARCH)
		}
		return ecsBinaryARM64, nil
	case "amd64":
		if len(ecsBinaryAMD64) == 0 {
			return nil, fmt.Errorf("AMD64 二进制文件未嵌入 (GOOS=%s, GOARCH=%s)", runtime.GOOS, runtime.GOARCH)
		}
		return ecsBinaryAMD64, nil
	default:
		return nil, fmt.Errorf("不支持的架构: %s (GOOS=%s)", runtime.GOARCH, runtime.GOOS)
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
