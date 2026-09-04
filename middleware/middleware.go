package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"mybbs/config"
	SH "mybbs/statehandler"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			SH.Error(c, http.StatusUnauthorized, "未认证", nil)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			SH.Error(c, http.StatusUnauthorized, "认证格式错误", nil)
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(config.CFG.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			SH.Error(c, http.StatusUnauthorized, "令牌无效或已过期", err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			SH.Error(c, http.StatusUnauthorized, "令牌解析失败", err)
			return
		}

		c.Set("user_id", uint(claims["user_id"].(float64)))
		c.Set("username", claims["username"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			SH.Error(c, http.StatusForbidden, "禁止访问", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
