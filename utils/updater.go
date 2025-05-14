package utils

import (
	"errors"
	"os"
	"time"
)

func update() error {
	if !IsAdmin() {
		err := ElevateSelf(os.Args)
		return err
	}
	today := TodayDate()
	status := runupdate(today)

	return status
}

func runupdate(today string) error {
	if !FileExists(CFG.ProfilePath) {
		Log("Profile文件不存在")
		Log("Profile file does not exist")
		return errors.New("Profile文件不存在")
	}

	Log("ReSet Trial Date:" + today)
	Log("重新设置试用日期:" + today)
	if err := backupFile(CFG.ProfilePath); err != nil {
		Log("备份出错...")
		Log("Backup failed...", err)
		return err
	}
	Log("备份Profile文件...")
	Log("Backing up Profile file...")

	config, err := readProfile(CFG.ProfilePath)
	if err != nil {
		Log("解析出错...")
		Log("Parse failed...", err)
		return err
	}
	Log("解析Profile文件...")
	Log("Parsing Profile file...")

	config["_iD"] = today
	//改写PreVersion 失效
	//config["version"] = "9.99.9"
	//禁用自动更新
	config["enableAutoUpdate"] = false
	//禁止发送匿名使用消息
	config["send_usage_info"] = false

	if err := writeProfile(CFG.ProfilePath, config); err != nil {
		Log("写入出错...")
		Log("Write failed...", err)
		restoreErr := restoreBackup(CFG.ProfilePath)
		if restoreErr != nil {
			Log("恢复失败，请手动检查 profile.data.bak")
			Log("Recovery failed, please manually check profile.data.bak")
		} else {
			Log("已成功恢复 profile.data.bak")
			Log("Successfully restored profile.data.bak")
		}
		return err
	}

	if err := updateRegistry(today, CFG.SID); err != nil {
		Log("写入注册表出错...")
		Log("Failed to write to registry...", err)
		return err
	}

	Log("已同步注册表...")
	Log("Registry synchronized...")
	return nil
}

func RunUpdaterLoop() {
	for {
		err := update()
		if err != nil {
			Log("更新失败:", err)
			Log("Update failed:", err)
		}
		Log("更新完成")
		Log("Update completed")
		time.Sleep(1 * time.Hour)
	}
}

func RunInteractive() {
	err := update()
	if err != nil {
		Log("更新失败:", err)
		Log("Update failed:", err)
	}
}
