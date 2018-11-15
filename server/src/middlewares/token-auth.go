package middlewares

import (
	"errors"

	"github.com/maple-leaf/happydoc-server/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Auth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("user")
		token := c.GetHeader("x-Token")
		user := models.User{Name: username}
		x := db.Where(user).First(&user)
		c.Header("X-Request-Id", "123")
		err := errors.New("auth error")
		if x.Error == nil {
			err = nil
		}
		if err != nil {
			c.JSON(403, gin.H{
				"message": "token invalid",
				"t":       token,
			})
		} else {
			c.JSON(200, gin.H{
				"message": token,
			})
		}
		c.Writer.Header().Set("X-Request-Id", "456")
		c.Next()
	}
}
