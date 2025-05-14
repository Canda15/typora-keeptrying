package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
)

func runasadmin(args []string) error {
	verb, err := syscall.UTF16PtrFromString("runas")
	if err != nil {
		return err
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err := syscall.UTF16PtrFromString(exePath)
	if err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	dir, err := syscall.UTF16PtrFromString(cwd)
	if err != nil {
		return err
	}

	// 拼接参数（忽略第一个参数 程序路径）
	paramsStr := strings.Join(args[1:], " ")
	params, err := syscall.UTF16PtrFromString(paramsStr)
	if err != nil {
		return err
	}

	r, _, _ := ShellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(exe)),
		uintptr(unsafe.Pointer(params)),
		uintptr(unsafe.Pointer(dir)),
		1,
	)
	if r <= 32 {
		return errors.New("ShellExecute failed: %v")
	}
	Log("已使用管理员权限重新运行：" + exePath + " " + paramsStr)
	Log("Restarted with admin privileges: " + exePath + " " + paramsStr)
	return nil
}

func updateRegistry(today string, sid string) error {
	regKey := fmt.Sprintf(`%s\Software\Typora`, sid)
	Log(regKey)
	//"S-1-5-21-2495842453-42734561234-229025492-1006\\Software\\Typora"
	k, _, err := registry.CreateKey(registry.USERS, regKey, registry.SET_VALUE)
	if err != nil {
		Log("无法打开或创建注册表项: " + err.Error())
		Log("Failed to open or create registry key: " + err.Error())
		return errors.New("无法打开或创建注册表项: %w")
	}
	defer k.Close()

	err = k.SetStringValue("IDate", today)
	if err != nil {
		Log("无法设置注册表值: " + err.Error())
		Log("Failed to set registry value: " + err.Error())
		return errors.New("无法设置注册表值: %w")
	}
	return nil
}

var (
	shell32      = syscall.NewLazyDLL("shell32.dll")
	ShellExecute = shell32.NewProc("ShellExecuteW")
)

func ConvertSIDToString(sid *windows.SID) (string, error) {
	var strPtr *uint16
	err := windows.ConvertSidToStringSid(sid, &strPtr)
	if err != nil {
		return "", err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(strPtr)))

	return windows.UTF16PtrToString(strPtr), nil
}

func GetLoggedInUserSID() (string, error) {
	var token windows.Token
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	if err != nil {
		Log("无法打开进程令牌: " + err.Error())
		Log("Failed to open process token: " + err.Error())
		return "", errors.New("无法打开进程令牌: %w")
	}
	defer token.Close()

	user, err := token.GetTokenUser()
	if err != nil {
		Log("无法获取令牌用户: " + err.Error())
		Log("Failed to get token user: " + err.Error())
		return "", errors.New("无法获取令牌用户: %w")
	}

	return ConvertSIDToString(user.User.Sid)
}

func IsAdmin() bool {
	k, err := registry.OpenKey(syscall.HKEY_LOCAL_MACHINE, "SOFTWARE", syscall.KEY_WRITE)
	if err == nil {
		k.Close()
		return true
	}
	return false
}

func TodayDate() string {
	return time.Now().Format("1/2/2006")
}

func ElevateSelf(args []string) error {
	Log("需要管理员权限，尝试重新以管理员运行...")
	Log("Admin privileges required, attempting to restart with admin rights...")
	err := runasadmin(args)
	if err != nil {
		Log("管理员权限提升失败:" + err.Error())
		Log("Failed to elevate privileges:" + err.Error())
		return err
	}
	return nil
}

func IsService() (bool, error) {
	isService, err := svc.IsWindowsService()
	return isService, err
}
