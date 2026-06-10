package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"kunkun-go/pkg/cosclient"

	"github.com/spf13/viper"
)

var cosUploader *cosclient.Client

// InitCOSUpload 从环境变量与 Viper 初始化 COS（应在 config.InitConfig 之后调用）。
func InitCOSUpload() {
	secretID := os.Getenv("COS_SECRET_ID")
	secretKey := os.Getenv("COS_SECRET_KEY")
	region := viper.GetString("upload.cos.region")
	bucket := viper.GetString("upload.cos.bucket")
	publicBase := viper.GetString("upload.cos.public_base_url")

	c, err := cosclient.New(secretID, secretKey, region, bucket, publicBase)
	if err != nil {
		panic(fmt.Errorf("init cos upload: %w", err))
	}
	cosUploader = c
}

// COSUploader 返回已初始化的 COS 客户端（供测试或扩展用）。
func COSUploader() *cosclient.Client {
	return cosUploader
}

// UploadToCOS 统一上传入口：校验大小与后缀，写入 COS，返回可访问 URL。
func UploadToCOS(ctx context.Context, header *multipart.FileHeader) (string, error) {
	if cosUploader == nil {
		return "", fmt.Errorf("COS 未初始化")
	}
	maxBytes := viper.GetInt64("upload.max_bytes")
	if maxBytes <= 0 {
		maxBytes = 10 << 20
	}
	if header.Size > maxBytes {
		return "", fmt.Errorf("文件超过大小限制（最大 %d 字节）", maxBytes)
	}

	ext := strings.ToLower(strings.TrimPrefix(path.Ext(header.Filename), "."))
	if ext == "" {
		return "", fmt.Errorf("无法识别文件扩展名")
	}
	if !allowedExt(ext) {
		return "", fmt.Errorf("不允许的文件类型: .%s", ext)
	}

	f, err := header.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	const sniffLen = 512
	head := make([]byte, sniffLen)
	n, err := io.ReadFull(f, head)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", err
	}
	mimeType := http.DetectContentType(head[:n])
	if mimeType == "application/octet-stream" {
		if mt := mime.TypeByExtension("." + ext); mt != "" {
			mimeType = strings.Split(mt, ";")[0]
		}
	}

	body := io.MultiReader(bytes.NewReader(head[:n]), f)
	objectKey := buildObjectKey(ext)

	return cosUploader.Upload(ctx, objectKey, body, header.Size, mimeType)
}

func allowedExt(ext string) bool {
	list := viper.GetStringSlice("upload.allowed_ext")
	if len(list) == 0 {
		list = []string{"jpg", "jpeg", "png", "gif", "webp", "pdf", "zip"}
	}
	for _, a := range list {
		if strings.EqualFold(strings.TrimPrefix(a, "."), ext) {
			return true
		}
	}
	return false
}

func buildObjectKey(ext string) string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	id := hex.EncodeToString(b[:])
	date := time.Now().UTC().Format("2006/01/02")
	return fmt.Sprintf("uploads/%s/%s.%s", date, id, ext)
}
