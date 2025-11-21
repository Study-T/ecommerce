package service

import (
	"archive/zip"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type UploadService struct {
	BasePath string
}

func NewUploadService(basePath string) *UploadService {
	return &UploadService{
		BasePath: basePath,
	}
}

// UploadFolder 上传文件夹（支持空文件夹）
func (s *UploadService) UploadFolder(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("无法打开上传文件: %w", err)
	}
	defer src.Close()

	// 创建目标目录
	folderName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	targetPath := filepath.Join(s.BasePath, folderName)

	// 检查文件类型
	if filepath.Ext(file.Filename) == ".zip" {
		return s.extractZipWithEmptyFolders(src, targetPath)
	}

	return "", fmt.Errorf("不支持的文件格式，请上传 ZIP 文件")
}

// extractZipWithEmptyFolders 解压ZIP文件并保留空文件夹
func (s *UploadService) extractZipWithEmptyFolders(src multipart.File, targetPath string) (string, error) {
	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-*.zip")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 复制上传的文件到临时文件
	if _, err := io.Copy(tempFile, src); err != nil {
		return "", fmt.Errorf("复制文件失败: %w", err)
	}

	// 打开ZIP文件
	zipReader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("打开ZIP文件失败: %w", err)
	}
	defer zipReader.Close()

	// 创建根目录
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return "", fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 解压所有文件和文件夹
	for _, file := range zipReader.File {
		filePath := filepath.Join(targetPath, file.Name)

		// 检查是否为目录（包括空目录）
		if file.FileInfo().IsDir() {
			// 创建空目录
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return "", fmt.Errorf("创建目录失败 %s: %w", file.Name, err)
			}
			// 创建 .gitkeep 文件以便Git跟踪空目录
			gitkeepPath := filepath.Join(filePath, ".gitkeep")
			if _, err := os.Stat(gitkeepPath); os.IsNotExist(err) {
				if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
					// 忽略 .gitkeep 创建失败的错误
					fmt.Printf("警告: 无法创建 .gitkeep 文件: %v\n", err)
				}
			}
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return "", fmt.Errorf("创建父目录失败: %w", err)
		}

		// 创建文件
		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return "", fmt.Errorf("创建文件失败 %s: %w", file.Name, err)
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return "", fmt.Errorf("打开ZIP文件内容失败: %w", err)
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()

		if err != nil {
			return "", fmt.Errorf("写入文件失败 %s: %w", file.Name, err)
		}
	}

	return targetPath, nil
}

// UploadFile 上传单个文件
func (s *UploadService) UploadFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("无法打开上传文件: %w", err)
	}
	defer src.Close()

	// 创建目标文件路径
	targetPath := filepath.Join(s.BasePath, file.Filename)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建目标文件
	dst, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	return targetPath, nil
}
