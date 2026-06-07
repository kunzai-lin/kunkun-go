package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Claims 访问令牌载荷（解析后从 [Claims.UserID] 取登录用户 id）。
type Claims struct {
	UserID uint `json:"user_id"`
	jwtlib.RegisteredClaims
}

var (
	// ErrMissingSecret 未配置 jwt.secret。
	ErrMissingSecret = errors.New("jwt.secret is not configured")
	// ErrWeakSecret 密钥过短，HS256 下建议至少 16 字节。
	ErrWeakSecret = errors.New("jwt.secret is too short (use at least 16 characters)")
)

func signingKey() ([]byte, error) {
	s := viper.GetString("jwt.secret")
	if s == "" {
		return nil, ErrMissingSecret
	}
	if len(s) < 16 {
		return nil, ErrWeakSecret
	}
	return []byte(s), nil
}

func tokenTTL() time.Duration {
	h := viper.GetInt("jwt.expire_hours")
	if h <= 0 {
		h = 72
	}
	return time.Duration(h) * time.Hour
}

// GenerateToken 为指定用户 id 签发 HS256 JWT。
func GenerateToken(userID uint) (string, error) {
	key, err := signingKey()
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(tokenTTL())),
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
		},
	}
	t := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return t.SignedString(key)
}

// ParseToken 校验并解析 JWT，返回载荷；签名错误、过期或格式无效时返回 error。
func ParseToken(tokenString string) (*Claims, error) {
	key, err := signingKey()
	if err != nil {
		return nil, err
	}
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// ParseBearer 从 "Bearer <token>" 或裸 token 字符串中解析 JWT。
func ParseBearer(authorizationHeader string) (*Claims, error) {
	s := strings.TrimSpace(authorizationHeader)
	if s == "" {
		return nil, errors.New("missing authorization")
	}
	if strings.HasPrefix(strings.ToLower(s), "bearer ") {
		s = strings.TrimSpace(s[7:])
	}
	return ParseToken(s)
}
