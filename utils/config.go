// utils/config.go
package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	ProfilePath    string `json:"profile_path"` // C:\Users\Administrator\AppData\Roaming\Typora
	InstallDirPath string `json:"install_path"` // C:\Typora_KeepTrying
	SID            string `json:"security_id"`  // S-1-5-21-2495842453-42734561234-229025492-1006
}

var DefaultInstall = `C:\Typora_KeepTrying`
var CFG *Config = &Config{}

const ServiceName = "TyporaFreeService"

func InitConfig(server bool) error {
	var err error
	CFG, err = LoadConfig()
	if err != nil {
		if server {
			return err
		}
		if CFG.InstallDirPath == "" {
			CFG.InstallDirPath = DefaultInstall
		}
		if CFG.ProfilePath == "" {
			CFG.ProfilePath, err = GetProfilePath()
			if err != nil {
				return err
			}
		}
		if CFG.SID == "" {
			CFG.SID, err = GetLoggedInUserSID()
			if err != nil {
				return err
			}
		}
		SaveConfig(CFG)
	}
	if server {
		return nil
	}
	Log("当前配置: \n", "安装路径：", CFG.InstallDirPath, " \nTypora数据: ", CFG.ProfilePath)
	Log("Current configuration: \n", "Installation path: ", CFG.InstallDirPath, " \nTypora data: ", CFG.ProfilePath)
	return nil
}

func GetConfigPath() string {
	var path string
	if CFG.InstallDirPath != "" {
		path := filepath.Join(CFG.InstallDirPath, "config.json")
		if FileExists(path) {
			Log("发现配置文件：", path)
			Log("Found configuration file:", path)
			return path
		}
	}
	path = filepath.Join(GetExecutableDir(), "config.json")
	if FileExists(path) {
		Log("发现配置文件：", path)
		Log("Found configuration file:", path)
		return path
	}

	path = filepath.Join(DefaultInstall, "config.json")
	if FileExists(path) {
		Log("发现配置文件：", path)
		Log("Found configuration file:", path)
		return path
	}

	path = ""
	Log("未发现配置文件", path)
	Log("No configuration file found", path)
	return path
}

func SaveConfig(cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(cfg.InstallDirPath), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	Log("保存配置", cfg.ProfilePath)
	Log("Saving configuration", cfg.ProfilePath)
	return os.WriteFile(filepath.Join(cfg.InstallDirPath, "config.json"), data, 0644)
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(GetConfigPath())
	if err != nil {
		Log("读取出错")
		Log("Read error")
		return &Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		Log("解析出错")
		Log("Parse error")
		return &Config{}, err
	}

	if cfg.ProfilePath == "" {
		Log("profile_path 为空")
		Log("profile_path is empty")
		return &cfg, errors.New("profile_path 为空")
	}

	return &cfg, nil
}

// 不带保存的加载配置
func CFGLoad() error {
	var err error
	if CFG, err = LoadConfig(); err != nil {
		Log("加载配置失败")
		Log("Failed to load configuration")
		if CFG.InstallDirPath == "" {
			CFG.InstallDirPath = DefaultInstall
		}
		if CFG.ProfilePath == "" {
			CFG.ProfilePath, err = GetProfilePath()
			if err != nil {
				return err
			}
		}
		Log("载入配置:", CFG)
		Log("Loaded configuration:", CFG)
		Log("当前配置: \n", "安装路径：", CFG.InstallDirPath, " \nTypora数据: ", CFG.ProfilePath)
		Log("Current configuration: \n", "Installation path: ", CFG.InstallDirPath, " \nTypora data: ", CFG.ProfilePath)
		return err
	}
	return nil
}
