package middleware

import (
	"kunkun-go/internal/repository"
	"kunkun-go/pkg/jwt"
	"kunkun-go/pkg/response"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CtxUserIDKey 通过 [JWTAuth] 写入 gin.Context，handler 里用 c.Get(CtxUserIDKey) 取当前用户 id。
const CtxUserIDKey = "userID"

// JWTAuth 校验 Authorization: Bearer <token>，通过后把用户 id 写入上下文。
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := jwt.ParseBearer(c.GetHeader("Authorization"))
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "未登录或 token 无效")
			c.Abort()
			return
		}

		// 检查用户是否在 Redis 中标记为已退出（须在 claims 有效之后访问 UserID）
		if repository.RDB != nil {
			logoutKey := "logout:" + strconv.Itoa(int(claims.UserID))
			exists, _ := repository.RDB.Exists(c.Request.Context(), logoutKey).Result()
			if exists == 1 {
				response.Error(c, http.StatusUnauthorized, "用户已退出登录")
				c.Abort()
				return
			}
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Next()
	}
}
