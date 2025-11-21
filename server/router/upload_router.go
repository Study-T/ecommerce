package router

import (
	"ecommerce/service"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UploadRouter struct {
	uploadService *service.UploadService
}

func NewUploadRouter(uploadService *service.UploadService) *UploadRouter {
	return &UploadRouter{
		uploadService: uploadService,
	}
}

func (r *UploadRouter) RegisterRoutes(router *gin.Engine) {
	uploadGroup := router.Group("/api/upload")
	{
		uploadGroup.POST("/folder", r.UploadFolder)
		uploadGroup.POST("/file", r.UploadFile)
	}
}

// UploadFolder 处理文件夹上传（支持空文件夹）
func (r *UploadRouter) UploadFolder(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未找到上传文件",
			"error":   err.Error(),
		})
		return
	}

	// 验证文件格式
	if filepath.Ext(file.Filename) != ".zip" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "只支持 ZIP 格式的文件夹上传",
		})
		return
	}

	targetPath, err := r.uploadService.UploadFolder(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件夹上传成功（包括空文件夹）",
		"data": gin.H{
			"path":     targetPath,
			"filename": file.Filename,
			"size":     file.Size,
		},
	})
}

// UploadFile 处理单个文件上传
func (r *UploadRouter) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未找到上传文件",
			"error":   err.Error(),
		})
		return
	}

	targetPath, err := r.uploadService.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "上传失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
		"data": gin.H{
			"path":     targetPath,
			"filename": file.Filename,
			"size":     file.Size,
		},
	})
}
