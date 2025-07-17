# mysqldeploy

# 功能
可在linux系统上安装 mysql5.*/8.0/8.4 版本的 MySQL，可以初始化多实例 MySQL。

# 限制
仅在 centos7 环境进行测试，后期可能支持更多系统。
此程序仅安装mysql大版本下的最终版本。

# 用法
交互式安装，程序会从官网下载安装包，并自动化安装。
```
# 下载安装命令
wget https://github.com/hhkens/mysqldeploy/releases/download/2.5.06/mysqldeploy_2.5.06.zip

# 解压
unzip mysqldeploy_2.5.06.zip

# 部署 MySQL
[root@localhost]# ./mysqldeploy 
基础检查已完成

请选择操作类型: 
        1) 安装新 MySQL
        2) 初始化 MySQL 多实例
        0) 退出程序

请输入类型编号[1|2|0] (默认: 1): 1
请选择版本: 
    1) MySQL-5.5
    2) MySQL-5.6
    3) MySQL-5.7
    4) MySQL-8.0
    5) MySQL-8.4
    0) 退出程序

请输入版本编号[1|2|3|4|5|0]: 
```
# 登录
可使用快捷方式登录
```
# 命令行输入 3306_mysql_login 输入 root 密码即可登录
[root@localhost]# 3306_mysql_login 
Enter password: # 输入密码即可
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 6
Server version: 5.7.44-log MySQL Community Server (GPL)

Copyright (c) 2000, 2023, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> 
mysql> 
mysql> 
mysql> 
```
快捷登录脚本位置: /usr/local/bin/3306_mysql_login 

# 目录结构
安装目录: /usr/local/mysqlxx  # xx: 版本号
数据目录: 可随意修改
端口: 可随意修改
配置文件: 数据目录下的 etc/my.cnf

# 安装过程
mysql 5.7 安装过程
```
[root@localhost]# ./mysqldeploy   # 执行命令
基础检查已完成

请选择操作类型: 
        1) 安装新 MySQL
        2) 初始化 MySQL 多实例
        0) 退出程序

请输入类型编号[1|2|0] (默认: 1):  # 选择类型 1 为安装  2 为初始化新实例

请选择版本: 
    1) MySQL-5.5
    2) MySQL-5.6
    3) MySQL-5.7
    4) MySQL-8.0
    5) MySQL-8.4
    0) 退出程序

请输入版本编号[1|2|3|4|5|0]: 3  # 选择安装mysql5.7
请选择端口[默认: 3306]: 
请选择数据目录[默认: /usr/local/mysql57/mysql3306 ]: 
请选择内存大小(单位: GB)[默认: 1GB]: 
请配置root密码[默认: 123456]: 
MySQL 配置信息如下:
使用端口: 3306
数据目录: /usr/local/mysql57/mysql3306
内存大小: 1 GB
root密码: 123456
请确认配置后继续 [Y/N]: y
2025-07-17 18:26:45 [info ] 安装包下载中···
2025-07-17 18:26:45 [info ] 安装依赖包···
2025-07-17 18:26:49 [info ] 解压安装包···
2025-07-17 18:26:49 [info ] 开始解压 mysql57.tar.gz 到 /usr/local/mysql57···
2025-07-17 18:27:04 [info ] 解压完成: /usr/local/mysql57
2025-07-17 18:27:04 [info ] mysql用户已存在
2025-07-17 18:27:04 [info ] 配置文件已生成
2025-07-17 18:27:04 [info ] 正在初始化···
2025-07-17 18:27:08 [info ] MySQL初始化成功
2025-07-17 18:27:08 [info ] 警告: 无法立即生效环境变量,请手动执行 source /etc/profile.d/mysql3306.sh 或重新登录
2025-07-17 18:27:08 [info ] MySQL服务启动成功
2025-07-17 18:27:08 [info ] MySQL开机自启设置成功
2025-07-17 18:27:08 [info ] 端口3306尚未就绪,等待3s后重试 (1/20)...
2025-07-17 18:27:11 [info ] MySQL服务端口3306已成功监听
2025-07-17 18:27:11 [info ] MySQL快捷登录脚本已创建: /usr/local/bin/3306_mysql_login
2025-07-17 18:27:11 [info ] 使用方式: 3306_mysql_login 
2025-07-17 18:27:11 [info ]  
2025-07-17 18:27:11 [info ] ==========================================
2025-07-17 18:27:11 [info ] ====安装目录: /usr/local/mysql57
2025-07-17 18:27:11 [info ] ====数据目录: /usr/local/mysql57/mysql3306
2025-07-17 18:27:11 [info ] ====root密码: 123456
2025-07-17 18:27:11 [info ] 便捷登录命令: 3306_mysql_login
2025-07-17 18:27:11 [info ] ==========================================
2025-07-17 18:27:11 [info ] MySQL安装完成!

```
