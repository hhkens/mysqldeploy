package utils

import (
	"bytes"
	"os/exec"
	"strconv"
	"text/template"
	"time"

	"fmt"
	"mysqldeploy/internal/conf"
	"mysqldeploy/internal/logs"
	"os"
	"path/filepath"

	"golang.org/x/exp/rand"
)

type MySQLTmpl struct {
	User                 string
	DataDir              string
	BaseDir              string
	TmpDir               string
	Socket               string
	PIDFile              string
	ServerID             int
	InnoDBBufferPoolSize string
	Port                 string
	ErrLog               string
	SlowLog              string
	BinLog               string
	BinLogIndex          string
	RelayLog             string
	RelayLogIndex        string
	GeneralLog           string
	XSocket              string
}

type MySQLStartTmpl struct {
	MySQLDCommand string
	MySQLCnfPath  string
}

// Decompress 解压.tar.gz或.tar.xz到指定路径
// filePath: 要解压的文件完整路径
// targetPath: 完整目标目录路径（如 "/opt/mysql57"）
func Decompress(filePath string, targetPath string) error {
	// 创建目标目录（确保存在）
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 构造最简tar命令
	cmd := exec.Command("tar", "xf", filePath, "-C", targetPath, "--strip-components=1")

	// 捕获命令输出
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// 执行解压命令
	logs.Info("开始解压 %s 到 %s···", filepath.Base(filePath), targetPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("解压失败: %w 错误输出: %s", err, stderr.String())
	}

	logs.Info("解压完成: %s", targetPath)
	return nil
}

// 创建目录
func CreateDir(paths []string) error {
	for _, path := range paths {
		path := filepath.Clean(path)
		if path == "" {
			return fmt.Errorf("无效的目录: 空字符串")
		}

		if _, err := os.Stat(path); err == nil {
			logs.Info("目录[%s]已存在", path)
			continue // 已存在则跳过
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("创建目录失败[%s]:%v", path, err)
		}
		if info, err := os.Stat(path); err != nil {
			if !info.IsDir() {
				return fmt.Errorf("路径存在但不是目录: %s", path)
			}
		}
	}
	return nil
}

// 创建配置文件
func ConfFile(cfg *conf.Config, mysqlUser, tmplPath string) error {
	rand.Seed(uint64(time.Now().UnixNano())) // 初始化随机种子
	num := rand.Intn(1000000)                // 生成0-999999的随机数

	conf := MySQLTmpl{
		User:                 mysqlUser,
		DataDir:              filepath.Join(cfg.DataDir, "data"),
		BaseDir:              cfg.InstallDir,
		TmpDir:               filepath.Join(cfg.DataDir, "tmp"),
		Socket:               filepath.Join(cfg.DataDir, "socket", "mysql.sock"),
		PIDFile:              filepath.Join(cfg.DataDir, "pid", "mysql.pid"),
		ServerID:             num,
		InnoDBBufferPoolSize: cfg.BuferrSize,
		Port:                 cfg.Port,
		ErrLog:               filepath.Join(cfg.DataDir, "log", "mysql-err.log"),
		SlowLog:              filepath.Join(cfg.DataDir, "log", "mysql-slow.log"),
		BinLog:               filepath.Join(cfg.DataDir, "binlog", "mysql-bin"),
		BinLogIndex:          filepath.Join(cfg.DataDir, "binlog", "mysql-bin.index"),
		RelayLog:             filepath.Join(cfg.DataDir, "relaylog", "relay-bin"),
		RelayLogIndex:        filepath.Join(cfg.DataDir, "relaylog", "relay-bin.index"),
		GeneralLog:           filepath.Join(cfg.DataDir, "log", "mysql-general.log"),
		XSocket:              filepath.Join(cfg.DataDir, "socket", "mysqlx.sock"),
	}

	// tmpl, err := template.ParseFiles("templates/my57.cnf.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	confFile := filepath.Join(cfg.DataDir, "etc", "my.cnf")

	file, err := os.Create(confFile)
	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, conf); err != nil {
		return fmt.Errorf("err: %v", err)
	}

	return nil
}

// 配置启动文件并启动 systemd
func StartFile(cfg *conf.Config, tmplPath string) error {
	// 1. 准备模板数据
	c := MySQLStartTmpl{
		MySQLDCommand: filepath.Join(cfg.InstallDir, "bin", "mysqld"),
		MySQLCnfPath:  filepath.Join(cfg.DataDir, "etc", "my.cnf"), // 修正了ect->etc
	}

	// 2. 验证关键路径是否存在
	if _, err := os.Stat(c.MySQLDCommand); os.IsNotExist(err) {
		return fmt.Errorf("mysqld二进制文件不存在: %s", c.MySQLDCommand)
	}

	// 3. 解析模板
	// tmpl, err := template.ParseFiles("templates/mysql57.service.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("模板解析失败: %v", err)
	}

	// 4. 创建systemd服务文件
	serviceFile := filepath.Join("/usr/lib/systemd/system", fmt.Sprintf("mysqld%s.service", cfg.Port))

	// 使用600权限创建文件，确保安全
	file, err := os.OpenFile(serviceFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("创建服务文件失败: %v", err)
	}
	defer file.Close()

	// 5. 执行模板写入
	if err := tmpl.Execute(file, c); err != nil {
		return fmt.Errorf("模板写入失败: %v", err)
	}

	// 6. 重载systemd
	if output, err := exec.Command("systemctl", "daemon-reload").CombinedOutput(); err != nil {
		return fmt.Errorf("daemon-reload失败: %v\n输出: %s", err, string(output))
	}

	// 7. 启动服务
	serviceName := fmt.Sprintf("mysqld%s.service", cfg.Port)
	if output, err := exec.Command("systemctl", "start", serviceName).CombinedOutput(); err != nil {
		return fmt.Errorf("服务启动失败: %v\n输出: %s", err, string(output))
	}
	logs.Info("MySQL服务启动成功")

	// 8. 设置开机自启
	if output, err := exec.Command("systemctl", "enable", serviceName).CombinedOutput(); err != nil {
		logs.Error("开机自启设置失败(非致命错误): %v\n输出: %s", err, string(output))
		// 这里不返回错误，因为服务已成功启动
	} else {
		logs.Info("MySQL开机自启设置成功")
	}
	// 9. 转换端口号为整数
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return fmt.Errorf("端口号转换失败: %v", err)
	}

	// 10. 循环检查端口是否启动成功，最多检查20次，每次间隔3秒
	maxRetries := 20
	retryInterval := time.Second * 3
	for i := 0; i < maxRetries; i++ {
		if CheckPortInUse(port) {
			logs.Info("MySQL服务端口%d已成功监听", port)
			return nil
		}

		if i < maxRetries-1 {
			logs.Info("端口%d尚未就绪,等待%v后重试 (%d/%d)...",
				port, retryInterval, i+1, maxRetries)
			time.Sleep(retryInterval)
		}
	}

	return nil
}

// 设置快捷登录
func Login(cfg *conf.Config) error {
	// 1. 验证关键路径
	mysqlCommand := filepath.Join(cfg.InstallDir, "bin", "mysql")
	if _, err := os.Stat(mysqlCommand); os.IsNotExist(err) {
		return fmt.Errorf("mysql客户端不存在于: %s", mysqlCommand)
	}

	socketFile := filepath.Join(cfg.DataDir, "socket", "mysql.sock")
	if _, err := os.Stat(socketFile); os.IsNotExist(err) {
		return fmt.Errorf("MySQL socket文件不存在于: %s", socketFile)
	}

	// 2. 构建脚本内容（更安全的密码处理）
	scriptContent := fmt.Sprintf(`#!/bin/bash
%s -S %s -p
`, mysqlCommand, socketFile)

	// 3. 创建目标目录（如果需要）
	targetDir := "/usr/local/bin"

	// 4. 创建脚本文件（限制权限）
	scriptPath := filepath.Join(targetDir, cfg.Port+"_mysql_login")
	file, err := os.OpenFile(scriptPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return fmt.Errorf("创建脚本文件失败: %v", err)
	}
	defer file.Close()

	// 5. 写入内容
	if _, err := file.WriteString(scriptContent); err != nil {
		return fmt.Errorf("写入脚本内容失败: %v", err)
	}

	logs.Info("MySQL快捷登录脚本已创建: %s", scriptPath)
	logs.Info("使用方式: %s ", cfg.Port+"_mysql_login")
	return nil
}
