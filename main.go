package main

import (
	"fmt"
	"mysqldeploy/internal/conf"
	"mysqldeploy/internal/initialize"
	"mysqldeploy/internal/install"
	"mysqldeploy/internal/logs"
	"mysqldeploy/internal/utils"
	"os"
)

func init() {

	if err := utils.BaseCheck(); err != nil {
		logs.Error(err.Error())
		os.Exit(2)
	}
	fmt.Println("基础检查已完成")

	// 选择运行类型 安装/初始化
	if err := utils.SelectType(); err != nil {
		logs.Error(err.Error())
		os.Exit(2)
	}
}

func main() {

	// 选择安装或初始化新实例
	switch conf.MysqlType {
	case "Install":
		// 初始化默认参数
		if err := mysqlInstall(); err != nil {
			return
		}
		logs.Info("MySQL安装完成!")
	case "Initialization":
		if err := initialization(); err != nil {
			return
		}
	}

	fmt.Println("")
}

func mysqlInstall() error {
	// 选择版本
	if err := utils.Select(); err != nil {
		logs.Error(err.Error())
		os.Exit(2)
	}

	cfg := conf.InitConfig()

	// 选择配置并设置
	if err := utils.SelectConf(cfg); err != nil {
		logs.Error("错误信息: %v", err)
		os.Exit(2)
	}

	// 检查配置
	if err := utils.Base2Check(cfg); err != nil {
		logs.Error("错误信息: %v", err)
		os.Exit(2)
	}

	switch conf.MysqlVersion {
	case "mysql55":
		if err := install.MySQL55(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	case "mysql56":
		if err := install.MySQL56(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	case "mysql57":
		if err := install.MySQL57(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	case "mysql80":
		if err := install.MySQL80(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	case "mysql84":
		if err := install.MySQL84(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	default:
		return fmt.Errorf("版本选择有误: %v", conf.MysqlVersion)
	}
	return nil
}

// 初始化 MySQL
func initialization() error {
	// 选择版本
	if err := utils.Select(); err != nil {
		logs.Error(err.Error())
		os.Exit(2)
	}

	cfg := conf.InitConfig()

	// 选择配置并设置
	if err := utils.SelectConf(cfg); err != nil {
		logs.Error("错误信息: %v", err)
		os.Exit(2)
	}

	// 检查配置
	if err := utils.InitializeBase2Check(cfg); err != nil {
		logs.Error("错误信息: %v", err)
		os.Exit(2)
	}

	switch conf.MysqlVersion {
	case "mysql55":
		if err := initialize.MySQL55(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	case "mysql56":
		if err := initialize.MySQL56(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	case "mysql57":
		if err := initialize.MySQL57(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	case "mysql80":
		if err := initialize.MySQL80(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	case "mysql84":
		if err := initialize.MySQL84(cfg); err != nil {
			logs.Error(err.Error())
			return err
		}
	default:
		return fmt.Errorf("版本选择有误: %v", conf.MysqlVersion)
	}
	return nil
}
