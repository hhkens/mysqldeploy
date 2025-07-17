package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Download(url, destDir, fileName, md5Key string) error {

	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 发起HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器返回错误: %s", resp.Status)
	}

	// 创建目标文件
	filePath := filepath.Join(destDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 初始化进度读取器
	sct := &ProgressReader{
		Reader:   resp.Body,
		Total:    resp.ContentLength,
		FileName: fileName,
	}

	// 执行下载
	if _, err := io.Copy(file, sct); err != nil {
		// 下载失败时删除不完整文件
		os.Remove(filePath)
		return fmt.Errorf("下载写入失败: %w", err)
	}

	// 验证文件完整性
	if fi, err := file.Stat(); err == nil && resp.ContentLength > 0 {
		if fi.Size() != resp.ContentLength {
			os.Remove(filePath)
			return fmt.Errorf("文件大小不匹配，可能下载不完整")
		}
	}
	if err := DiffMD5(filePath, md5Key); err == nil {
		return err
	}
	return nil
}

// 确保 ProgressReader 实现 io.Reader 接口
var _ io.Reader = (*ProgressReader)(nil)

type ProgressReader struct {
	Reader   io.Reader
	Total    int64
	Current  int64
	FileName string
	lastTime time.Time
}

func (r *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.Current += int64(n)

	// 每秒最多更新一次显示
	if time.Since(r.lastTime) > time.Second {
		// 转换为MB计算
		currentMB := float64(r.Current) / (1024 * 1024)
		totalMB := float64(r.Total) / (1024 * 1024)
		percent := float64(r.Current) / float64(r.Total) * 100

		fmt.Printf("\r%s: %.2f%% [%.2f/%.2f MB]",
			r.FileName,
			percent,
			currentMB,
			totalMB)
		r.lastTime = time.Now()
	}
	return
}

func DiffMD5(filePath string, md5Key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := md5.New()
	buf := make([]byte, 16*1024*1024) // 建议4MB-16MB

	if _, err := io.CopyBuffer(hash, file, buf); err != nil {
		return err
	}

	pmd5 := hex.EncodeToString(hash.Sum(nil))

	if pmd5 != md5Key {
		return fmt.Errorf("md5值不同,下载失败")
	}

	return nil
}
