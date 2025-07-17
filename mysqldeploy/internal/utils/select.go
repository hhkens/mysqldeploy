package utils

import (
	"bufio"
	"fmt"
	"mysqldeploy/internal/conf"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 选择安装类型 安装/初始化
func SelectType() error {
	text := `
请选择操作类型: 
	1) 安装新 MySQL
	2) 初始化 MySQL 多实例
	0) 退出程序
`
	for {
		fmt.Println(text)
		fmt.Printf("请输入类型编号[1|2|0] (默认: 1): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取输入失败: %v", err)
		}
		input = strings.TrimSpace(input)
		input = strings.TrimSpace(input)

		if input == "" {
			input = "1"
		}
		// fmt.Printf("input: %v\n", input)

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			fmt.Println("无效选择, 请输入 0-2 之间的数字")
			continue
		}
		// fmt.Printf("choice: %v\n", choice)
		if choice < 0 || choice > 5 {
			fmt.Println("无效选择, 请输入 0-2 之间的数字")
			continue
		}
		switch choice {
		case 1:
			conf.MysqlType = "Install"
		case 2:
			conf.MysqlType = "Initialization"
		case 0:
			return fmt.Errorf("用户选择退出程序")
		default:
			return fmt.Errorf("选择无效")
		}
		// fmt.Println(conf.MysqlType)
		return nil
	}
}

// 选择安装的版本
func Select() error {
	text := `
请选择版本: 
    1) MySQL-5.5
    2) MySQL-5.6
    3) MySQL-5.7
    4) MySQL-8.0
    5) MySQL-8.4
    0) 退出程序
`
	for {
		fmt.Println(text)
		fmt.Printf("请输入版本编号[1|2|3|4|5|0]: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取输入失败: %v", err)
		}
		input = strings.TrimSpace(input)
		input = strings.TrimSpace(input)
		// fmt.Printf("input: %v\n", input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			fmt.Println("无效选择, 请输入 0-5 之间的数字")
			continue
		}
		// fmt.Printf("choice: %v\n", choice)
		if choice < 0 || choice > 5 {
			fmt.Println("无效选择, 请输入 0-5 之间的数字")
			continue
		}
		switch choice {
		case 1:
			conf.MysqlVersion = "mysql55"
		case 2:
			conf.MysqlVersion = "mysql56"
		case 3:
			conf.MysqlVersion = "mysql57"
		case 4:
			conf.MysqlVersion = "mysql80"
		case 5:
			conf.MysqlVersion = "mysql84"
		case 0:
			return fmt.Errorf("用户选择退出程序")
		default:
			return fmt.Errorf("选择无效")
		}
		// fmt.Println(conf.MysqlVersion)
		return nil
	}
}

// 用户输入配置信息
func SelectConf(cfg *conf.Config) error {
	// fmt.Printf("conf.MysqlVersion: %v\n", conf.MysqlVersion)

	fmt.Printf("请选择端口[默认: %v]: ", cfg.Port)
	port, err := GetUserInput(cfg.Port)
	if err != nil {
		return err
	}
	cfg.Port = port

	dataDir := filepath.Join(cfg.InstallDir, "mysql"+cfg.Port)
	fmt.Printf("请选择数据目录[默认: %v ]: ", dataDir)
	input, err := GetUserInput(dataDir)
	if err != nil {
		return err
	}
	if input != dataDir {
		input = filepath.Join(input, "mysql"+cfg.Port)
	}
	cfg.DataDir = input

	fmt.Printf("请选择内存大小(单位: GB)[默认: %vGB]: ", cfg.BuferrSize)

	memSize, err := GetUserInput(cfg.BuferrSize)
	if err != nil {
		return err
	}
	cfg.BuferrSize = memSize

	fmt.Printf("请配置root密码[默认: %v]: ", cfg.Password)
	pwd, err := GetUserInput(cfg.Password)
	if err != nil {
		return err
	}
	cfg.Password = pwd

	fmt.Println("MySQL 配置信息如下:")
	fmt.Printf("使用端口: %v\n", cfg.Port)
	fmt.Printf("数据目录: %v\n", cfg.DataDir)
	fmt.Printf("内存大小: %v GB\n", cfg.BuferrSize)
	fmt.Printf("root密码: %v\n", cfg.Password)

	fmt.Print("请确认配置后继续 [Y/N]: ")
	status, err := GetUserInput("")
	if err != nil {
		return err
	}
	status = strings.ToUpper(status)
	if status != "Y" {
		return fmt.Errorf("用户选择结束运行")
	}
	return nil
}
