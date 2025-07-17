package initialize

import (
	"fmt"
	"mysqldeploy/internal/conf"
	"mysqldeploy/internal/logs"
	"mysqldeploy/internal/utils"
)

func MySQL80(cfg *conf.Config) error {

	mysqlUser := "mysql"

	// 创建目录
	paths := cfg.GetFullPaths()

	if err := utils.CreateDir(paths); err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// 生成配置文件
	if err := utils.ConfFile(cfg, mysqlUser, "templates/my80.cnf.tmpl"); err != nil {
		return err
	}
	logs.Info("配置文件已生成")

	// 授权用户
	if err := utils.ChownUser(mysqlUser, cfg.InstallDir); err != nil {
		return err
	}
	if err := utils.ChownUser(mysqlUser, cfg.DataDir); err != nil {
		return err
	}

	// 初始化
	logs.Info("正在初始化···")
	if err := utils.Initialize(cfg, mysqlUser); err != nil {
		return err
	}

	// 配置启动脚本并启动
	if err := utils.StartFile(cfg, "templates/mysql80.service.tmpl"); err != nil {
		return err
	}
	// 快捷登录
	if err := utils.Login(cfg); err != nil {
		return err
	}

	// 设置密码
	if err := utils.SetPassword(cfg); err != nil {
		return err
	}

	if err := utils.Prompt(cfg); err != nil {
		return err
	}

	return nil
}
