package utils

import (
	"bufio"
	"fmt"
	"mysqldeploy/internal/logs"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// GetUserInput 获取用户输入
// prompt: 提示信息
// defaultValue: 默认值(可选)
// validator: 输入验证函数(可选)
func GetUserInput(defaultValue string) (string, error) {
	// 读取用户输入
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("读取输入失败: %v", err)
	}
	// 清理输入(去除前后空格和换行符)
	input = strings.TrimSpace(input)
	// fmt.Printf("input: %v\n", input)
	// 如果输入为空且提供了默认值，则使用默认值
	if input == "" && defaultValue != "" {
		input = defaultValue
	}
	return input, nil
}

// CreateSystemUser 创建一个无家目录且不能登录的系统用户
// username: 用户名
// shell: 设置为 /usr/sbin/nologin 或 /bin/false 禁止登录
func CreateUser(username, shell string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("仅支持Linux")
	}

	if shell == "" {
		shell = "/usr/sbin/nologin"
	}

	// 检查用户是否已存在
	if UserExists(username) {
		logs.Info("%v用户已存在", username)
		return nil
	}

	// 使用 useradd 创建用户（无家目录、无登录权限）
	cmd := exec.Command("sudo", "useradd",
		"--no-create-home", // 不创建家目录
		"--shell", shell,   // 设置不能登录的shell
		username,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("创建用户失败: %v, output: %s", err, string(output))
	}

	return nil
}

// 检查用户是否存在
func UserExists(username string) bool {
	cmd := exec.Command("id", "-u", username)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// 授权目录
func ChownUser(username, path string) error {
	userGroup := username + ":" + username
	cmd := exec.Command("chown", "-R", userGroup, path)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("chown failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}
