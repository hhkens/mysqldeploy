package conf

import (
	"path/filepath"
)

var (
	MysqlType    string // 全新安装/添加新实例
	MysqlVersion string // MySQL 版本选则
)

var MySQLConfig Config

const (
	Mysql55DownloadURL = "https://downloads.mysql.com/archives/get/p/23/file/mysql-5.5.62-linux-glibc2.12-x86_64.tar.gz"
	Mysql55DownloadMD5 = "2cb52d5ca4eef4cd37783ab5cb3622f3"

	Mysql56DownloadURL = "https://downloads.mysql.com/archives/get/p/23/file/mysql-5.6.51-linux-glibc2.12-x86_64.tar.gz"
	Mysql56DownloadMD5 = "b4a58f228dc5e2d579eb9ce59e566bd5"

	Mysql57DownloadURL = "https://downloads.mysql.com/archives/get/p/23/file/mysql-5.7.44-linux-glibc2.12-x86_64.tar.gz"
	Mysql57DownloadMD5 = "d7c8436bbf456e9a4398011a0c52bc40"

	Mysql80DownloadURL = "https://downloads.mysql.com/archives/get/p/23/file/mysql-8.0.41-linux-glibc2.17-x86_64.tar.xz"
	Mysql80DownloadMD5 = "736d3800cf5e8504eddf0701c42c2ad5"

	Mysql84DownloadURL = "https://downloads.mysql.com/archives/get/p/23/file/mysql-8.4.4-linux-glibc2.17-x86_64.tar.xz"
	Mysql84DownloadMD5 = "683db607b34406af7fc9080690776df1"
)

type Config struct {
	Port       string
	InstallDir string
	DataDir    string
	SubDirs    []string
	BuferrSize string
	Password   string
}

func InitConfig() *Config {
	cfg := &Config{Port: "3306"}
	cfg.InstallDir = filepath.Join("/usr/local/", MysqlVersion)
	cfg.DataDir = filepath.Join(cfg.InstallDir)
	cfg.SubDirs = []string{"binlog", "data", "etc", "lock", "log", "pid", "socket", "tmp", "script", "relaylog"}
	cfg.BuferrSize = "1"
	cfg.Password = "123456"

	return cfg
}

func (c *Config) GetFullPaths() []string {
	var paths []string
	for _, path := range c.SubDirs {
		paths = append(paths, filepath.Join(c.DataDir, path))

	}
	return paths
}
