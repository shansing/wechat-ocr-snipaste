//go:build cgo
// +build cgo

package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//const TEMP_DIR_NAME = "wechat-ocr"

func main() {
	exePath, err := os.Executable()
	if err != nil {
		printDebug("[Error] get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	stdinFlag := flag.Bool("stdin", false, "must set")
	stdoutFlag := flag.Bool("stdout", false, "must set")
	ocrBinArg := flag.String("ocrBin", filepath.Join(exeDir, "ocrBin\\wxocr.dll"), "full path (including file name) of wxocr.dll or wechatocr.exe")
	wechatDirArg := flag.String("wechatDir", filepath.Join(exeDir, "wechat\\wechatDir"), "WeChat path (usually with a version in path)")
	args := os.Args[1:]
	for i, arg := range args {
		if arg == "stdin" || arg == "stdout" {
			args[i] = "-" + arg
		}
	}
	err = flag.CommandLine.Parse(args)
	if err != nil {
		printDebug("[Error] bad args %v", args)
		os.Exit(1)
	}
	if *stdinFlag != true || *stdoutFlag != true {
		printDebug("[Error] \"stdin stdout\" must be set. This program is for Snipaste. %v", args)
		os.Exit(1)
	}
	printDebug("[Info] ocrBin: %s\n", *ocrBinArg)
	exists, err := pathExists(*ocrBinArg)
	if exists != true {
		printDebug("[Error] ocrBin seems not to exist")
		os.Exit(1)
	}
	if err != nil {
		printDebug("[Error] ocrBin access denied %w", err)
		os.Exit(1)
	}
	printDebug("[Info] wechatDir: %s\n", *wechatDirArg)
	exists, err = pathExists(*wechatDirArg)
	if exists != true {
		printDebug("[Error] wechatDir seems not to exist")
		os.Exit(1)
	}
	if err != nil {
		printDebug("[Error] wechatDir access denied %w", err)
		os.Exit(1)
	}

	//tempDir := filepath.Join(os.TempDir(), TEMP_DIR_NAME)
	//if err := os.MkdirAll(tempDir, 0o755); err != nil {
	//	printDebug("[Error] mkdir %s: %w", tempDir, err)
	//	os.Exit(1)
	//}

	uuid, err := randHex(16)
	if err != nil {
		printDebug("[Error] generate uuid: %w", err)
		os.Exit(1)
	}

	//filePath := filepath.Join(tempDir, uuid+".png")
	filePath := filepath.Join(os.TempDir(), uuid+".png")
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		printDebug("[Error] create file %s: %w", filePath, err)
		os.Exit(1)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(filePath)
	}()
	if _, err := io.Copy(f, os.Stdin); err != nil {
		printDebug("[Error] write stdin to %s: %w", filePath, err)
		os.Exit(1)
	}
	if err := f.Close(); err != nil {
		printDebug("[Error] close file %s: %w", filePath, err)
		os.Exit(1)
	}

	dllDir := filepath.Join(exeDir, "wcocr.dll")
	printDebug("[Info] dllDir: %s\n", dllDir)

	//printDebug("[Info] 1")
	wechatOCR, err := NewWechatOCR(dllDir)
	if err != nil {
		printDebug("[Error] Failed to create WechatOCR instance: %v\n", err)
		os.Exit(1)
	}
	//printDebug("[Info] 2")
	result := OcrCustom(wechatOCR, *ocrBinArg, *wechatDirArg, filePath)
	if result == nil || result.Errcode != 0 {
		printDebug("[Error] err code is not 0 %v", result)
		os.Exit(1)
	}
	//printDebug("[Info] 3")
	text := extractTexts(result)
	_, _, _ = wechatOCR.stop_ocr.Call()

	_, err = io.WriteString(os.Stdout, text)
	if err != nil {
		printDebug("[Error] write to stdout: %s", text)
		os.Exit(1)
	}
}

func printDebug(format string, a ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, a)
}

func extractTexts(result *Result) string {
	var builder strings.Builder
	for _, ocr := range result.OcrResponse {
		builder.WriteString(ocr.Text)
		builder.WriteByte('\n')
	}
	return builder.String()
}

func randHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "DEFALUT", err
	}
	return hex.EncodeToString(b), nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		// 文件或路径存在
		return true, nil
	}
	if os.IsNotExist(err) {
		// 文件或路径不存在
		return false, nil
	}
	// 发生了别的错误，比如权限问题，返回错误
	return false, err
}
