package middleware

import (
	"kunkun-go/pkg/jwt"
	"kunkun-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// CtxUserIDKey 通过 [JWTAuth] 写入 gin.Context，handler 里用 c.Get(CtxUserIDKey) 取当前用户 id。
const CtxUserIDKey = "userID"

// JWTAuth 校验 Authorization: Bearer <token>，通过后把用户 id 写入上下文。
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.ParseBearer(c.GetHeader("Authorization"))
		if err != nil {
			response.Error(c, 401, "未登录或 token 无效")
			c.Abort()
			return
		}
		c.Set(CtxUserIDKey, claims.UserID)
		c.Next()
	}
}
