package main

import (
	"fmt"
	"os"
	"strings"
	"typora-keeptrying/serve"
	"typora-keeptrying/utils"
)

func main() {
	args := os.Args
	if len(args) > 1 {
		command := strings.ToLower(args[1])
		switch command {
		case "install":
			if !utils.IsAdmin() {
				utils.ElevateSelf(args)
				return
			}
			serve.ExeInstall()
		case "del":
			if !utils.IsAdmin() {
				utils.ElevateSelf(args)
				return
			}
			if err := utils.CFGLoad(); err != nil {
				return
			}
			serve.Uninstall()
		case "start":
			if !utils.IsAdmin() {
				utils.ElevateSelf(args)
				return
			}
			serve.ControlService("start")
		case "stop":
			if !utils.IsAdmin() {
				utils.ElevateSelf(args)
				return
			}
			serve.ControlService("stop")
		case "restart":
			if !utils.IsAdmin() {
				utils.ElevateSelf(args)
				return
			}
			serve.ControlService("stop")
			serve.ControlService("start")
		default:
			fmt.Println("未知命令: ", command)
			fmt.Println("Unknown command:", command)

			fmt.Println("可用命令: install | del | start | stop | restart")
			fmt.Println("Available commands: install | del | start | stop | restart")
		}
		return
	}

	isService, _ := utils.IsService()
	utils.InitConfig(isService)
	if isService {
		serve.RunService()
	} else {
		utils.RunInteractive()
	}
	fmt.Println("按任意键退出...  (Press any key to exit...) ")
	fmt.Scanln()
}
