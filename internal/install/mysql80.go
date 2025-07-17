package install

import (
	"fmt"
	"mysqldeploy/internal/conf"
	"mysqldeploy/internal/logs"
	"mysqldeploy/internal/utils"
	"os"
	"path/filepath"
)

func MySQL80(cfg *conf.Config) error {
	// 下载路径
	destDir := "/tmp"
	// 解压包名
	fileName := "mysql80.tar.xz"

	mysqlUser := "mysql"

	// 压缩包+路径
	filePath := filepath.Join(destDir, fileName)
	// 解压目录+路径
	DecompressPath := filepath.Join(cfg.InstallDir)

	// 下载
	logs.Info("安装包下载中···")

	_, err := os.Stat(filePath)
	if err == nil { // 如果文件存在
		// 对比 MD5值
		if err := utils.DiffMD5(filePath, conf.Mysql80DownloadMD5); err != nil {
			err := utils.Download(conf.Mysql80DownloadURL, destDir, fileName, conf.Mysql80DownloadMD5)
			if err != nil {
				return fmt.Errorf("安装包下载失败: %v", err)
			}
		}
	} else {
		err = utils.Download(conf.Mysql80DownloadURL, destDir, fileName, conf.Mysql80DownloadMD5)
		if err != nil {
			return fmt.Errorf("安装包下载失败: %v", err)
		}
	}
	logs.Info("安装依赖包···")
	if err := utils.InstallDependencies(); err != nil {
		return fmt.Errorf("依赖安装失败: %v", err)
	}

	// 解压
	logs.Info("解压安装包···")
	if err := utils.Decompress(filePath, DecompressPath); err != nil {
		return fmt.Errorf("解压安装包失败: %v", err)
	}

	// 创建目录
	paths := cfg.GetFullPaths()

	if err := utils.CreateDir(paths); err != nil {
		return fmt.Errorf("error: %v", err)
	}
	// 创建用户
	if err := utils.CreateUser(mysqlUser, ""); err != nil {
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

	// 环境变量
	if err := utils.EnvVariables(cfg); err != nil {
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
