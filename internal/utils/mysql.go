package utils

import (
	"bytes"
	"database/sql"
	"fmt"
	"mysqldeploy/internal/conf"
	"mysqldeploy/internal/logs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 初始化 MySQL
func Initialize(cfg *conf.Config, mysqlUser string) error {
	mysqldPath := filepath.Join(cfg.InstallDir, "bin", "mysqld")
	baseDir := filepath.Join(cfg.InstallDir)
	dataDir := filepath.Join(cfg.DataDir, "data")

	// 确保目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %v", err)
	}

	// 正确构建命令行参数（每个参数单独传递）
	args := []string{
		"--initialize-insecure",
		"--user=" + mysqlUser,
		"--basedir=" + baseDir,
		"--datadir=" + dataDir,
		"--lower_case_table_names=1",
		"--innodb_data_file_path=ibdata1:1024M:autoextend",
	}

	cmd := exec.Command(mysqldPath, args...)

	// 捕获命令输出（有助于调试）
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logs.Info("初始化失败: %v", err)
		logs.Info("命令输出: %s", out.String())
		logs.Info("错误输出: %s", stderr.String())
		return fmt.Errorf("初始化失败: %v, 错误输出: %s", err, stderr.String())
	}

	logs.Info("MySQL初始化成功")
	return nil
}

// 初始化 MySQL
func Initialize56(cfg *conf.Config, mysqlUser string) error {
	mysqldPath := filepath.Join(cfg.InstallDir, "scripts", "mysql_install_db")
	baseDir := filepath.Join(cfg.InstallDir)
	dataDir := filepath.Join(cfg.DataDir, "data")

	// 确保目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("创建数据目录失败: %v", err)
	}

	// 正确构建命令行参数（每个参数单独传递）
	args := []string{
		"--user=" + mysqlUser,
		"--basedir=" + baseDir,
		"--datadir=" + dataDir,
		"--lower_case_table_names=1",
		"--innodb_data_file_path=ibdata1:1024M:autoextend",
	}

	cmd := exec.Command(mysqldPath, args...)

	// 捕获命令输出（有助于调试）
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logs.Info("初始化失败: %v", err)
		logs.Info("命令输出: %s", out.String())
		logs.Info("错误输出: %s", stderr.String())
		return fmt.Errorf("初始化失败: %v, 错误输出: %s", err, stderr.String())
	}

	logs.Info("MySQL初始化成功")
	return nil
}

// 安装依赖包
func InstallDependencies() error {
	var cmd *exec.Cmd

	// 根据操作系统选择命令
	switch runtime.GOOS {
	case "linux":
		// 检测是否是基于RPM的系统（CentOS/RHEL）
		if _, err := exec.LookPath("yum"); err == nil {
			cmd = exec.Command("sudo", "yum", "install", "-y",
				"libaio-devel",
				"libaio",
				"perl-Data-Dumper",
				"openssl-devel",
				"jemalloc")
		} else if _, err := exec.LookPath("apt-get"); err == nil {
			// Debian/Ubuntu 系统
			cmd = exec.Command("sudo", "apt-get", "install", "-y",
				"libaio-dev",
				"perl-base",
				"libdata-dumper-simple-perl",
				"libssl-dev",
				"libjemalloc-dev")
		} else {
			return fmt.Errorf("不支持的Linux发行版")
		}
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	// 执行并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("安装失败: %v\n输出: %s", err, string(output))
	}

	return nil
}

// 配置环境变量
func EnvVariables(cfg *conf.Config) error {
	mysqlBin := filepath.Join(cfg.InstallDir, "bin")

	// 确保目录存在
	profileDir := filepath.Join("/etc", "profile.d")

	// 构造文件路径
	filepath := filepath.Join(profileDir, fmt.Sprintf("mysql%s.sh", cfg.Port))

	// 构造文件内容
	content := fmt.Sprintf("export PATH=%s:$PATH\n", mysqlBin)

	// 以只写模式创建文件，如果存在则截断
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 写入内容
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	// 使配置立即生效（可选）
	cmd := exec.Command("source", filepath)
	if err := cmd.Run(); err != nil {
		// 这里可能会失败，因为source是shell内置命令
		// 可以记录日志但不作为错误返回
		logs.Info("警告: 无法立即生效环境变量,请手动执行 source %s 或重新登录", filepath)
	}

	return nil
}

// SetPassword 设置MySQL root用户密码
// cfg: 包含MySQL配置信息的结构体
func SetPassword(cfg *conf.Config) error {
	// 1. 验证输入参数
	if cfg.Port == "" {
		return fmt.Errorf("MySQL端口号不能为空")
	}
	if cfg.Password == "" {
		return fmt.Errorf("新密码不能为空")
	}

	// 2. 构造DSN连接字符串
	sockPath := filepath.Join(cfg.DataDir, "socket", "mysql.sock")
	dsn := fmt.Sprintf("root@unix(%s)/?timeout=5s", sockPath)

	// fmt.Printf("dsn: %v\n", dsn)

	// 3. 连接MySQL服务器
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("MySQL连接失败: %v (请确认已导入_ \"github.com/go-sql-driver/mysql\")", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("mysql Ping 失败: %v", err)
	}

	// 4. 设置连接参数
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0) // 禁用空闲连接

	var alterUserSQL string
	// 6. 执行密码修改
	if conf.MysqlVersion == "mysql56" || conf.MysqlVersion == "mysql55" {
		alterUserSQL = fmt.Sprintf("SET PASSWORD FOR 'root'@'localhost' = PASSWORD('%s')", cfg.Password)
		if _, err := db.Exec(alterUserSQL); err != nil {
			return fmt.Errorf("修改密码失败: %v (请确认有足够权限)", err)
		}
	} else {
		alterUserSQL = fmt.Sprintf("ALTER USER 'root'@'localhost' IDENTIFIED BY '%s'", cfg.Password)
		if _, err := db.Exec(alterUserSQL); err != nil {
			return fmt.Errorf("修改密码失败: %v (请确认有足够权限)", err)
		}
	}

	// 8. 验证新密码是否生效
	if err := verifyNewPassword(cfg); err != nil {
		return fmt.Errorf("密码验证失败: %v (请手动确认密码是否已更改)", err)
	}

	return nil
}

// verifyNewPassword 验证新密码是否生效
func verifyNewPassword(cfg *conf.Config) error {

	sockPath := filepath.Join(cfg.DataDir, "socket", "mysql.sock")
	dsn := fmt.Sprintf("root:%s@unix(%s)/?timeout=5s", cfg.Password, sockPath)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Second)
	return db.Ping()
}

// 完成提示
func Prompt(cfg *conf.Config) error {
	logs.Info(" ")
	logs.Info("==========================================")
	logs.Info("====安装目录: %v", cfg.InstallDir)
	logs.Info("====数据目录: %v", cfg.DataDir)
	logs.Info("====root密码: %v", cfg.Password)
	logs.Info("便捷登录命令: %v", cfg.Port+"_mysql_login")
	logs.Info("==========================================")

	return nil
}
