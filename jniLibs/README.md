# JNI Libraries 目录说明

这个目录用于存放 Android APK 中的 native libraries（ECS 二进制文件）。

## 目录结构

```
jniLibs/
├── arm64-v8a/       # ARM64 架构 (64位) - 主要目标
│   └── libgoecs.so
└── x86_64/          # x86_64 架构 (64位) - 模拟器
    └── libgoecs.so
```

## 如何准备二进制文件

从 ECS 项目编译 Linux 二进制文件，然后复制并重命名为 `.so` 文件：

```bash
# 1. 编译 Linux 二进制文件（使用 goreleaser 参数）
cd /path/to/ecs

# ARM64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
  -ldflags="-s -w -X main.version=1.0.0 -X main.arch=arm64 -checklinkname=0" \
  -o goecs-linux-arm64 ./

# AMD64 (x86_64)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-s -w -X main.version=1.0.0 -X main.arch=amd64 -checklinkname=0" \
  -o goecs-linux-amd64 ./

# 2. 复制到 jniLibs 目录并重命名为 .so
cd /path/to/goecs
cp /path/to/ecs/goecs-linux-arm64 jniLibs/arm64-v8a/libgoecs.so
cp /path/to/ecs/goecs-linux-amd64 jniLibs/x86_64/libgoecs.so

# 3. 设置执行权限
chmod 755 jniLibs/*/libgoecs.so
```

## 工作原理

1. 在打包 APK 时，Fyne 会自动将 `jniLibs/` 目录中的文件打包进 APK
2. Android 系统在安装 APK 时，会将这些 `.so` 文件提取到应用的 `nativeLibraryDir`（通常是 `/data/app/<package>/lib/<abi>/`）
3. 应用运行时，通过 `embedding/embed_android.go` 从 `nativeLibraryDir` 读取文件路径
4. 使用 `exec.Command()` 直接执行该路径的二进制文件
5. 这种方式不需要 root 权限，也不会触发 SELinux 限制

## 注意事项

- **文件必须命名为 `libgoecs.so`**（或其他以 `lib` 开头、`.so` 结尾的名称）
- **使用 Linux 编译参数，不是 Android**：虽然目标是 Android，但使用 `GOOS=linux` 编译
- 必须放在正确的 ABI 目录下（`arm64-v8a` 或 `x86_64`）
- 只需要这两个架构就够用了（覆盖真机和模拟器）

## 快速命令（一键准备）

假设 ECS 项目在 `../ecs`，当前在 `goecs` 项目根目录：

```bash
# 编译并复制
cd ../ecs && \
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -checklinkname=0" -o goecs-linux-arm64 ./ && \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -checklinkname=0" -o goecs-linux-amd64 ./ && \
cd ../goecs && \
cp ../ecs/goecs-linux-arm64 jniLibs/arm64-v8a/libgoecs.so && \
cp ../ecs/goecs-linux-amd64 jniLibs/x86_64/libgoecs.so && \
chmod 755 jniLibs/*/libgoecs.so && \
ls -lh jniLibs/*/libgoecs.so
```

## 测试验证

### 1. 验证 APK 包含 .so 文件

```bash
unzip -l your-app.apk | grep "lib/"
# 应该看到：
# lib/arm64-v8a/libgoecs.so
# lib/x86_64/libgoecs.so
```

### 2. 在设备上验证安装位置

```bash
# 安装 APK
adb install your-app.apk

# 查看安装位置
adb shell pm path com.oneclickvirt.goecs

# 查看 native library 目录
adb shell ls -l /data/app/com.oneclickvirt.goecs-*/lib/arm64/

# 测试执行
adb shell
run-as com.oneclickvirt.goecs
cd /data/app/com.oneclickvirt.goecs-*/lib/arm64/
./libgoecs.so --help
```

### 3. 查看应用日志

```bash
# 实时查看日志
adb logcat | grep -E "(goecs|oneclickvirt)"

# 只看错误
adb logcat *:E | grep goecs
```

## 故障排除

如果看到 "fork/exec" 错误：

1. **检查文件是否存在**
   ```bash
   adb shell run-as com.oneclickvirt.goecs find /data/app -name "libgoecs.so"
   ```

2. **检查编译架构是否匹配设备**
   ```bash
   adb shell getprop ro.product.cpu.abi
   # arm64-v8a -> 需要 jniLibs/arm64-v8a/libgoecs.so
   # x86_64 -> 需要 jniLibs/x86_64/libgoecs.so
   ```

3. **验证文件不是空的**
   ```bash
   ls -lh jniLibs/*/libgoecs.so
   # 文件应该有合理的大小（几MB到几十MB）
   ```

4. **检查应用日志中的详细错误信息**
   - embedding 代码会输出所有检查过的路径
   - 查看 logcat 获取详细信息

