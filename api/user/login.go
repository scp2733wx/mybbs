package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"mybbs/config"
	db "mybbs/database"
	SH "mybbs/statehandler"
)

type LoginGET struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var LG LoginGET
		if err := c.ShouldBindJSON(&LG); err != nil {
			SH.Error(c, http.StatusBadRequest, "登录信息获取错误", err)
			return
		}

		var userTemp User
		if err := db.DB.Where("username = ?", LG.Username).First(&userTemp).Error; err != nil {
			SH.Error(c, http.StatusNotFound, "用户不存在", err)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(userTemp.PasswordHash), []byte(LG.Password)); err != nil {
			SH.Error(c, http.StatusUnauthorized, "密码错误", err)
			return
		}

		expireHour := 24
		if expireHour <= 0 {
			expireHour = 24
		}
		claims := jwt.MapClaims{
			"user_id":  userTemp.ID,
			"username": LG.Username,
			"role":     userTemp.Role,
			"exp":      time.Now().Add(time.Hour * time.Duration(expireHour)).Unix(),
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.CFG.JWT.Secret))
		if err != nil {
			SH.Error(c, http.StatusInternalServerError, "生成JWT令牌失败", err)
			return
		}

		lresponse := JWT_user{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   int64(expireHour * 3600),
			User: &User{
				ID:       userTemp.ID,
				Username: userTemp.Username,
				Name:     userTemp.Name,
				Role:     userTemp.Role,
			},
		}

		SH.Success(c, lresponse)
	}
}
