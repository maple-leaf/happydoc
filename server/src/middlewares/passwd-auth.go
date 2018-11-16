package middlewares

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/gin-contrib/sessions"
	"github.com/maple-leaf/happydoc-server/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func PasswdAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		if username != nil {
			return
		}

		username = c.PostForm("name")
		if username == "" {
			c.Redirect(301, "/login")
			return
		}
		passwd := c.PostForm("passwd")
		sum := sha256.Sum256([]byte(passwd))
		token := base64.StdEncoding.EncodeToString(sum[:])
		user := models.User{Name: username.(string)}
		x := db.Where(user).First(&user)
		if x.Error != nil || user.Token != token {
			c.Redirect(301, "/login")
			return
		}
		session.Set("username", username)
		session.Save()
		c.Next()
	}
}
