package cosclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// Client 封装腾讯云 COS 上传与可访问 URL 拼接。
type Client struct {
	region        string
	bucket        string
	publicBaseURL string
	inner         *cos.Client
}

// New 使用永久密钥创建 COS 客户端。bucket 为控制台完整桶名，形如 mybucket-1250000000。
func New(secretID, secretKey, region, bucket, publicBaseURL string) (*Client, error) {
	if secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("cosclient: COS_SECRET_ID / COS_SECRET_KEY 不能为空（请配置 .env）")
	}
	if region == "" || bucket == "" {
		return nil, fmt.Errorf("cosclient: region 与 bucket 不能为空（请在 config.yaml 的 upload.cos 中配置）")
	}
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	if err != nil {
		return nil, err
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return &Client{
		region:        region,
		bucket:        bucket,
		publicBaseURL: strings.TrimSpace(publicBaseURL),
		inner:         client,
	}, nil
}

// ObjectPublicURL 返回对象在浏览器中可直接访问的 URL（桶需公有读或 public_base_url 为 CDN 等可读地址）。
func (c *Client) ObjectPublicURL(objectKey string) string {
	objectKey = strings.Trim(objectKey, "/")
	if c.publicBaseURL != "" {
		base := strings.TrimSuffix(c.publicBaseURL, "/")
		return base + "/" + escapeKeySegments(objectKey)
	}
	base := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", c.bucket, c.region)
	return strings.TrimSuffix(base, "/") + "/" + escapeKeySegments(objectKey)
}

func escapeKeySegments(objectKey string) string {
	parts := strings.Split(objectKey, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}

// Upload 将对象写入 COS，返回可访问 URL（取决于桶权限与 public_base_url）。
func (c *Client) Upload(ctx context.Context, objectKey string, r io.Reader, contentLength int64, contentType string) (string, error) {
	objectKey = strings.TrimPrefix(objectKey, "/")
	opt := &cos.ObjectPutOptions{}
	if contentType != "" || contentLength > 0 {
		h := &cos.ObjectPutHeaderOptions{}
		if contentType != "" {
			h.ContentType = contentType
		}
		if contentLength > 0 {
			h.ContentLength = contentLength
		}
		opt.ObjectPutHeaderOptions = h
	}
	_, err := c.inner.Object.Put(ctx, objectKey, r, opt)
	if err != nil {
		return "", err
	}
	return c.ObjectPublicURL(objectKey), nil
}
