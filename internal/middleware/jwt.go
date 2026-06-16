package middleware

import (
	"echo-framework/internal/config"
	"echo-framework/pkg/security"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type JWTConfig struct {
	Secret         []byte
	Issuer         string
	SigningMethod  jwt.SigningMethod
	ExpirationTime time.Duration
}

type CustomClaims struct {
	UserID   int    `json:"userID"`
	UserName string `json:"userName"`
	jwt.RegisteredClaims
}

var (
	jwtConfig *JWTConfig
	mutex     sync.Mutex
)

// Jwt JWT中间件
func Jwt() echo.MiddlewareFunc {
	excludedPaths := map[string]bool{
		"/login":  true,
		"/public": true,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 检查排除路由
			if excludedPaths[c.Request().URL.Path] {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, echo.Map{
					"error": "授权头缺失",
				})
			}

			tokenString, err := parseBearerToken(authHeader)
			if err != nil {
				return c.JSON(401, echo.Map{
					"error": "无效的token格式",
				})
			}

			claims, err := validateJWT(tokenString)
			if err != nil {
				return handleJWTError(c, err)
			}

			// 存储claims到上下文
			c.Set("jwt_claims", claims)
			return next(c)
		}
	}
}

// GenerateToken 生成token
func GenerateToken(userID int, userName string) (string, error) {
	conf, _ := config.LoadConfig()
	jwtConf, err := loadJwtConfig(conf)
	if err != nil {
		return "", err
	}

	claims := CustomClaims{
		UserID:   userID,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtConf.ExpirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConf.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwtConf.SigningMethod, claims)
	tokenString, err := token.SignedString(jwtConf.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateJWT 验证并解析token
func validateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwtConfig.SigningMethod {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		return jwtConfig.Secret, nil
	})
	if err != nil {
		zap.L().Error("failed to validate token: %v", zap.Error(err))
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// handleJWTError 处理JWT错误
func handleJWTError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		ResetJWTConfig()
		return c.JSON(401, echo.Map{
			"error": "token过期",
		})
	case errors.Is(err, jwt.ErrTokenMalformed):
		return c.JSON(401, echo.Map{
			"error": "无效的token格式",
		})
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return c.JSON(401, echo.Map{
			"error": "无效的token签名",
		})
	default:
		return c.JSON(401, echo.Map{
			"error": "认证失败",
		})
	}
}

// parseBearerToken 解析Bearer token
func parseBearerToken(header string) (string, error) {
	const bearerPrefix = "Bearer "
	if len(header) < len(bearerPrefix) || !strings.HasPrefix(header, bearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return header[len(bearerPrefix):], nil
}

// loadJwtConfig 加载JWT配置
func loadJwtConfig(conf *config.Config) (*JWTConfig, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if jwtConfig != nil {
		return jwtConfig, nil
	}

	randSecret, err := security.CryptoRandSecret(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %v", err)
	}
	secret := []byte(randSecret)

	signingMethod := jwt.GetSigningMethod(conf.Jwt.SigningMethod)
	if signingMethod == nil {
		return nil, fmt.Errorf("invalid signing method: %s", conf.Jwt.SigningMethod)
	}

	expirationTime, err := time.ParseDuration(conf.Jwt.ExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration time format: %s", conf.Jwt.ExpirationTime)
	}

	jwtConfig = &JWTConfig{
		Secret:         secret,
		Issuer:         conf.Jwt.Issuer,
		SigningMethod:  signingMethod,
		ExpirationTime: expirationTime,
	}
	return jwtConfig, nil
}

// ResetJWTConfig 重置jwt配置
func ResetJWTConfig() {
	mutex.Lock()
	defer mutex.Unlock()
	jwtConfig = nil
}
