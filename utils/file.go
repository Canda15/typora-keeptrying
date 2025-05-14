package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func readProfile(path string) (map[string]interface{}, error) {
	rawHex, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := hex.DecodeString(string(rawHex))
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func writeProfile(path string, config map[string]interface{}) error {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	newHex := hex.EncodeToString(jsonBytes)
	return os.WriteFile(path, []byte(newHex), 0644)
}
func backupFile(path string) error {
	backupPath := path + ".bak"
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return os.WriteFile(backupPath, input, 0644)
}
func restoreBackup(path string) error {
	backupPath := path + ".bak"
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}
	return os.WriteFile(path, backupData, 0644)
}

func GetExecutableDir() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("无法获取可执行路径: %v", err)
		log.Fatalf("Failed to get executable path: %v", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		log.Fatalf("无法解析符号链接: %v", err)
		log.Fatalf("Failed to resolve symbolic link: %v", err)
	}
	return filepath.Dir(exePath)
}

func GetProfilePath() (string, error) {
	userDir := os.Getenv("APPDATA")
	if userDir == "" {
		Log("未能获取 APPDATA 环境变量")
		Log("Failed to get APPDATA environment variable")
		return "", errors.New("无法获取当前用户的 APPDATA 环境变量")
	}
	path := filepath.Join(userDir, "Typora", "profile.data")
	if !FileExists(path) {
		Log("找不到Profile.data 请手动检查")
		Log("Profile.data not found, please check manually")
		return "", errors.New("未找到 Profile.data 请手动检查")
	}
	Log("已找到Profile.data :", path)
	Log("Found Profile.data:", path)
	return path, nil
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 确保目标文件被覆盖
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 生成清理bat 等待程序占用解除后删除文件夹
func SelfDelete() error {
	batContent := `@echo off
:loop
tasklist | findstr /i "typora_free.exe" >nul
if not errorlevel 1 (
    timeout /t 1 >nul
    goto loop
)
rmdir /s /q "` + CFG.InstallDirPath + `"
del "%~f0"`
	batPath := filepath.Join(os.TempDir(), "cleanup.bat")
	err := os.WriteFile(batPath, []byte(batContent), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("cmd", "/C", "start", "", batPath)
	err = cmd.Start()
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func RemoveFilesWithExtensions(dirPath string, extensions []string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, ext := range extensions {
			if strings.HasSuffix(strings.ToLower(info.Name()), ext) {
				if err := os.Remove(path); err != nil {
					Log("删除文件 失败:", path, err)
					Log("Failed to delete file:", path, err)
				} else {
					Log("删除文件 成功:", path)
					Log("Successfully deleted file:", path)
				}
			}
		}
		return nil
	})
	return err
}

func OpenExploer(path string) {
	cmd := exec.Command("explorer", path)
	cmd.Start()
}
