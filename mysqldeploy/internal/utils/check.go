package utils

import (
	"bufio"
	"fmt"
	"io"
	"mysqldeploy/internal/conf"
	"mysqldeploy/internal/logs"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 基础检查
func BaseCheck() error {
	// 检查版本信息
	// 检查网络
	// 检查架构
	// 检查是否为 root 用户

	if !checkSysVersion() {
		fmt.Println("当前系统非 CentOS 7")
		fmt.Print("仅对 CentOS 7 进行验证是否继续?[Y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		input = strings.ToUpper(input)
		if input != "Y" {
			return fmt.Errorf("用户选择退出程序")
		}
		fmt.Println("用户选择继续执行···")
	}

	if !isNetworkAvailable() {
		logs.Error("当前无法访问 www.mysql.com")
		logs.Error("本程序依赖网络请确认网络状态后重试")
		return fmt.Errorf("连接互联网失败")
	}
	if !checkArch() {
		logs.Error("本程序仅支持x86_64架构")
		return fmt.Errorf("架构检查异常程序退出")
	}

	if os.Getegid() != 0 {
		return fmt.Errorf("请使用 root 用户运行此程序")
	}

	return nil
}

// 二次检查
func Base2Check(cfg *conf.Config) error {
	// 检查安装目录
	if !DirCheck(cfg.InstallDir) {
		return fmt.Errorf("安装目录: %v 不为空或是个文件,请处理后重试", cfg.InstallDir)
	}
	// 检查数据目录
	if !DirCheck(cfg.DataDir) {
		return fmt.Errorf("数据目录: %v 不为空或是个文件,请处理后重试", cfg.DataDir)
	}

	// 检查端口
	num, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return err
	}
	if CheckPortInUse(num) {
		return fmt.Errorf("%v 端口已存在", num)
	}

	return nil
}

// 二次检查
func InitializeBase2Check(cfg *conf.Config) error {
	// 检查安装目录
	if DirCheck(cfg.InstallDir) {
		return fmt.Errorf("安装目录: %v 为空,请处理后重试", cfg.InstallDir)
	}
	// 检查数据目录
	if !DirCheck(cfg.DataDir) {
		return fmt.Errorf("数据目录: %v 不为空或是个文件,请处理后重试", cfg.DataDir)
	}

	// 检查端口
	num, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return err
	}
	if CheckPortInUse(num) {
		return fmt.Errorf("%v 端口已存在", num)
	}

	return nil
}

// 检查架构
func checkArch() bool {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logs.Error("架构信息: %v", string(output))
		logs.Error(err.Error())
		return false
	}

	str := string(output)
	// fmt.Printf("str: %v\n", str)
	if !strings.Contains(str, "x86_64") {
		logs.Error("架构信息: %v", string(output))
		return false
	}
	return true
}

// 检查系统版本
func checkSysVersion() bool {
	// 检查系统版本
	f, err := os.Open("/etc/os-release")
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	defer f.Close()

	var sysID, versionID string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			sysID = strings.Trim(strings.TrimPrefix(line, "ID="), `"`)
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			versionID = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), `"`)
		}
	}
	// fmt.Printf("sysID: %v\n", sysID)
	// fmt.Printf("versionID: %v\n", versionID)

	return sysID == "centos" && versionID == "7"
}

// 检查网络
func isNetworkAvailable() bool {
	timeout := time.Duration(5 * time.Second)
	conn, err := net.DialTimeout("tcp", "www.mysql.com:80", timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// CheckPortInUse 检查指定端口是否被占用
// 参数 port: 要检查的端口号
// 返回 bool: true 表示端口已被占用，false 表示端口可用
func CheckPortInUse(port int) bool {
	// 将端口转换为字符串
	portStr := strconv.Itoa(port)

	// 尝试监听该端口
	listener, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		// 如果出错，说明端口可能已被占用
		return true
	}

	// 成功监听，说明端口可用，关闭监听器
	listener.Close()
	return false
}

// DirCheck 检查目录是否存在或是否为空
// 返回true表示目录不存在或目录为空
// 返回false表示目录存在且非空
func DirCheck(path string) bool {
	// 检查目录是否存在
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return true // 目录不存在
		}
		return false // 其他错误视为目录存在（保守处理）
	}

	// 检查是否是目录
	if !fileInfo.IsDir() {
		return false // 路径存在但不是目录
	}

	// 检查目录是否为空
	dir, err := os.Open(path)
	if err != nil {
		return false // 无法打开目录视为非空
	}
	defer dir.Close()

	// 读取最多2个文件（.和..不算）
	names, err := dir.Readdirnames(2)
	if err != nil && err != io.EOF {
		return false // 读取错误视为非空
	}

	// 过滤掉"."和".."
	var filtered []string
	for _, name := range names {
		if name != "." && name != ".." {
			filtered = append(filtered, name)
		}
	}

	return len(filtered) == 0
}
