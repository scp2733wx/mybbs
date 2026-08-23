package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	SH "mybbs/statehandler"
)

type RegisterGET struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var RG RegisterGET
		if err := c.ShouldBindJSON(&RG); err != nil {
			SH.Error(c, http.StatusInternalServerError, "注册信息获取错误："+err.Error())
			return
		}

		var userTemp User
		err := db.Where("username = ?", RG.Username).First(&userTemp).Error
		if err == nil {
			SH.Error(c, http.StatusConflict, "用户名重复")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(RG.Password), bcrypt.DefaultCost)
		if err != nil {
			SH.Error(c, http.StatusInternalServerError, "哈希加密失败："+err.Error())
			return
		}

		Newuser := User{
			Username:     RG.Username,
			Name:         RG.Name,
			PasswordHash: string(hashedPassword),
			Role:         RG.Role,
		}
		if err := db.Create(&Newuser).Error; err != nil {
			SH.Error(c, http.StatusInternalServerError, "注册用户数据错误："+err.Error())
			return
		}

		SH.JSON(c, http.StatusOK, 0, "注册成功", gin.H{
			"id":       Newuser.ID,
			"username": Newuser.Username,
			"name":     Newuser.Name,
			"role":     Newuser.Role,
		})
	}
}
