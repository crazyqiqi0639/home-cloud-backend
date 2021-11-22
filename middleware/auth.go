package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"home-cloud/models"
	"home-cloud/utils"
	"net/http"
)

func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username, ok := session.Get("user").(string)
		if !ok || len(username) == 0 {
			if c.Request.URL.Path != "/api/file/get_file" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": 1, "message": "You have not logged in!"})
			} else {
				c.String(http.StatusUnauthorized, "401 Unauthorized")
				c.Abort()
			}
			return
		}
		user, err := models.GetUserByUsername(username)
		if err != nil {
			if c.Request.URL.Path != "/api/file/get_file" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": 1, "message": "You have not logged in!"})
			} else {
				c.String(http.StatusUnauthorized, "401 Unauthorized")
				c.Abort()
			}
			return
		}
		utils.GetLogger().Info("User " + user.Username + " request comes")
		c.Set("user", user)
		c.Next()
	}
}

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.Value("user").(*models.User)
		if user.Status != 1 {
			utils.GetLogger().Warn("User " + user.Username + " try to access admin page, rejected")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": 1, "message": "Permission denied! "})
		} else {
			c.Next()
		}
	}
}
