package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/maple-leaf/happydoc-server/models"
	"github.com/maple-leaf/happydoc-server/services"
)

func JWT(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := c.GetHeader("x-Token")
		username := c.PostForm("account")
		user := models.User{
			Name: username,
		}
		q := db.Where(user).First(&user)
		if q.Error != nil {
			c.Status(403)
			c.Next()
			return
		}

		_, err := services.ValidateJWTRS(t, user.Token, services.AuthClaims{
			Username: username,
		})
		if err != nil {
			c.JSON(403, gin.H{
				"message": "token invalid",
				"t":       t,
			})
		} else {
			c.JSON(200, gin.H{
				"message": t,
			})
		}
		c.Next()
	}
}
