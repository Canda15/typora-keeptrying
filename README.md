这是一个无聊的适用于 Windows 平台 在 （1.10.8 上通过) 的无限试用 15 天脚本...........因为它真的快到期了所以......

它做了什么操作？

1. 尝试获取管理员权限
2. 记录当前用户 SID
3. 修改文件 C:\Users\UserName\AppData\Roaming\Typora\profile.data
4. 修改注册表 HKEY_CURRENT_USER\Software\Typora\IDate
5. 支持作为 Windows 服务自动运行重置程序（每 1h) ，提供基本的本地运行日志

它支持一些命令行参数

```
 install | del | start | stop | restart
```

你可以这样启动它

```bash
# go run . or ./typora-keeptrying.exe
# go run . install or ./typora-keeptrying.exe install
# go run . del or ./typora-keeptrying.exe del
# go run . restart or ./typora-keeptrying.exe restart
```

在 config.json 中 设置安装目录

```json
{
  "profile_path": "",
  "install_path": "C:\\Typora_KeepTrying"
}
```

如果运行目录中没有设置或不存在 `config.json` 则默认 安装到 `C:\\Typora_KeepTrying`
