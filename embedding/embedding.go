package embedding

// GetECSBinaryPath 提取并返回嵌入的 ECS 二进制文件路径
// 这个函数会根据不同平台调用相应的实现
func GetECSBinaryPath() (string, error) {
	return ExtractECSBinary()
}

// Cleanup 清理提取的二进制文件
func Cleanup(path string) {
	CleanupECSBinary(path)
}
