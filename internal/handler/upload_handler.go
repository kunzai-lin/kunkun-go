package handler

import (
	"net/http"

	"kunkun-go/internal/service"
	"kunkun-go/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// UploadFile 接收 multipart 字段 "file"，上传到腾讯云 COS，返回可访问 URL。
func UploadFile(c *gin.Context) {
	maxBytes := viper.GetInt64("upload.max_bytes")
	if maxBytes <= 0 {
		maxBytes = 10 << 20
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "请选择要上传的文件（字段名 file）")
		return
	}

	url, err := service.UploadToCOS(c.Request.Context(), file)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, gin.H{"url": url})
}
