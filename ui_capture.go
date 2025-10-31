package main

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
)

// captureOutputOnly 只捕获函数输出，不打印到终端（GUI专用）
func captureOutputOnly(f func()) string {
	// 保存旧的 stdout 和 stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// 创建管道
	stdoutPipeR, stdoutPipeW, err := os.Pipe()
	if err != nil {
		return "Error creating stdout pipe"
	}
	stderrPipeR, stderrPipeW, err := os.Pipe()
	if err != nil {
		stdoutPipeW.Close()
		stdoutPipeR.Close()
		return "Error creating stderr pipe"
	}

	// 替换标准输出和标准错误输出为管道写入端
	os.Stdout = stdoutPipeW
	os.Stderr = stderrPipeW

	// 恢复标准输出和标准错误输出
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		stdoutPipeW.Close()
		stderrPipeW.Close()
		stdoutPipeR.Close()
		stderrPipeR.Close()
	}()

	// 缓冲区 - 只捕获，不打印
	var stdoutBuf, stderrBuf bytes.Buffer

	// 并发读取 stdout 和 stderr
	done := make(chan struct{})
	go func() {
		io.Copy(&stdoutBuf, stdoutPipeR) // 只写入缓冲区，不打印到终端
		done <- struct{}{}
	}()
	go func() {
		io.Copy(&stderrBuf, stderrPipeR) // 只写入缓冲区，不打印到终端
		done <- struct{}{}
	}()

	// 执行函数
	f()

	// 关闭管道写入端，让管道读取端可以读取所有数据
	stdoutPipeW.Close()
	stderrPipeW.Close()

	// 等待两个 goroutine 完成
	<-done
	<-done

	// 返回捕获的输出字符串
	return stdoutBuf.String()
}

// printAndCaptureGUI 捕获函数输出的同时不打印到终端（GUI专用）
func printAndCaptureGUI(f func(), tempOutput, output string) string {
	tempOutput = captureOutputOnly(f)
	output += tempOutput
	return output
}

// appendResult 将文本追加到结果显示区域，并移除ANSI转义序列
func (ui *TestUI) appendResult(text string) {
	// 去除 ANSI 转义序列（颜色代码等）
	ansiRegex := regexp.MustCompile(`\x1B\[[0-9;]*[a-zA-Z]`)
	cleanText := ansiRegex.ReplaceAllString(text, "")

	currentText := ui.resultText.Text
	ui.resultText.SetText(currentText + cleanText)
	ui.resultText.CursorRow = len(strings.Split(ui.resultText.Text, "\n"))
}

// setResult 一次性设置所有结果
func (ui *TestUI) setResult(text string) {
	ui.resultText.SetText(text)
	ui.resultText.CursorRow = len(strings.Split(text, "\n"))
}
