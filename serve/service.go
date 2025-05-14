package serve

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
	"typora-keeptrying/utils"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

func ControlService(action string) {
	m, err := mgr.Connect()
	if err != nil {
		utils.Log("连接服务管理器失败:", err)
		utils.Log("Failed to connect to service manager:", err)
		return
	}
	defer m.Disconnect()

	s, err := m.OpenService(utils.ServiceName)
	if err != nil {
		utils.Log("打开服务失败:", err)
		utils.Log("Failed to open service:", err)
		return
	}
	defer s.Close()

	switch action {
	case "start":
		err = s.Start()
	case "stop":
		_, err = s.Control(svc.Stop)
	}

	if err != nil {
		utils.Log("服务操作失败 ["+action+"]:", err)
		utils.Log("Service operation failed ["+action+"]:", err)
	} else {
		utils.Log("服务操作成功 [" + action + "]")
		utils.Log("Service operation successful [" + action + "]")
	}
}

func RunService() {
	err := svc.Run(utils.ServiceName, &myService{})
	if err != nil {
		utils.Log("服务运行失败:", err)
		utils.Log("Service failed to run:", err)
	}
}

func ExeInstall() {
	exists, err := ServiceExists(utils.ServiceName)
	if err != nil {
		utils.Log("获取服务状态失败:", err)
		utils.Log("Failed to get service status:", err)
		return
	}
	if exists {
		utils.Log("检测到旧服务，清理中...")
		utils.Log("Detected old service, cleaning up...")
		cleanupOldService()
		cleanfiles()
	}
	if err := utils.InitConfig(false); err != nil {
		utils.Log("初始化配置失败:", err)
		utils.Log("Failed to initialize configuration:", err)
		return
	}

	//./typora-keeptrying.exe -> C:/Typora_KeepTrying/typora-keeptrying.exe
	exePath, err := os.Executable()
	if err != nil {
		utils.Log("获取执行路径失败:", err)
		utils.Log("Failed to get executable path:", err)
		return
	}

	destPath := filepath.Join(utils.CFG.InstallDirPath, filepath.Base(exePath))

	if err := os.MkdirAll(utils.CFG.InstallDirPath, 0755); err != nil {
		utils.Log("创建安装目录失败:", err)
		utils.Log("Failed to create installation directory:", err)
		return
	}

	if err := utils.CopyFile(exePath, destPath); err != nil {
		utils.Log("复制执行文件失败:", err)
		utils.Log("Failed to copy executable file:", err)
		return
	}

	err = registerService(destPath)
	if err != nil {
		utils.Log("服务注册失败:", err)
		utils.Log("Failed to register service:", err)
		return
	}

	utils.Log("服务注册完成")
	utils.Log("Service registration completed")
	utils.Log("安装路径:", utils.CFG.InstallDirPath)
	utils.Log("Installation path:", utils.CFG.InstallDirPath)
	utils.OpenExploer(utils.CFG.InstallDirPath)
}

// registerService 注册服务
func registerService(exePath string) error {
	m, err := mgr.Connect()
	if err != nil {
		return errors.New("无法连接到服务管理器: %v")
	}
	defer m.Disconnect()

	// 尝试打开现有的服务
	s, err := m.OpenService(utils.ServiceName)
	if err == nil {
		s.Close()
		utils.Log("服务已存在")
		utils.Log("Service already exists")
		return nil
	}

	// 创建新服务
	s, err = m.CreateService(utils.ServiceName, exePath, mgr.Config{
		DisplayName: "Typora Free Service",
		StartType:   mgr.StartAutomatic,
	})
	if err != nil {
		return errors.New("创建服务失败: %v")
	}
	defer s.Close()

	// 启动服务
	err = s.Start()
	if err != nil {
		return errors.New("启动服务失败: %v")
	}

	return nil
}

func cleanupOldService() error {
	m, err := mgr.Connect()
	if err != nil {
		utils.Log("无法连接服务管理器:", err)
		utils.Log("Failed to connect to service manager:", err)
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(utils.ServiceName)
	if err != nil {
		utils.Log(err)
		return nil
	}

	defer s.Close()

	status, err := s.Query()
	if err == nil && status.State == svc.Running {
		utils.Log("服务正在运行，尝试停止...")
		utils.Log("Service is running, attempting to stop...")
		_, err := s.Control(svc.Stop)
		if err != nil {
			utils.Log("停止服务失败:", err)
			utils.Log("Failed to stop service:", err)
			return err
		}
		// 等待服务停止
		for i := 0; i < 10; i++ {
			time.Sleep(500 * time.Millisecond)
			status, _ = s.Query()
			if status.State != svc.Running {
				break
			}
		}
		utils.Log("服务已停止")
		utils.Log("Service has been stopped")
	}

	// 删除服务
	err = s.Delete()
	if err != nil {
		utils.Log("删除旧服务失败:", err)
		utils.Log("Failed to delete old service:", err)
		return err
	} else {
		utils.Log("旧服务已删除")
		utils.Log("Old service has been deleted")
	}
	return nil
}

func ServiceExists(name string) (bool, error) {
	m, err := mgr.Connect()
	if err != nil {
		return false, err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		if strings.Contains(err.Error(), "The specified service does not exist") {
			return false, nil
		}
		return false, err
	}
	s.Close()
	return true, nil
}

func cleanfiles() error {
	if err := utils.RemoveFilesWithExtensions(utils.CFG.InstallDirPath, []string{".exe", ".json"}); err != nil {
		utils.Log("清理文件失败", err)
		utils.Log("Failed to clean files", err)
		return err
	}
	return nil
}

func Uninstall() error {
	cleanupOldService()
	// cleanfiles()
	utils.SelfDelete()
	return nil
}

type myService struct{}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	s <- svc.Status{State: svc.StartPending}
	go utils.RunUpdaterLoop()
	s <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
	for c := range r {
		switch c.Cmd {
		case svc.Interrogate:
			s <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			s <- svc.Status{State: svc.StopPending}
			return false, 0
		default:
			utils.Log("收到未知服务控制请求")
			utils.Log("Received unknown service control request")
		}
	}
	return false, 0
}
