package login

import (
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"meme/db"
	"meme/helpers"
	"net/http"
)

type Body struct {
	Name     string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Handler(c *gin.Context) {
	user := Body{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sql := "SELECT * FROM users WHERE username=?"
	queryUser := db.MemeDB.QueryRow(sql, user.Name)

	var (
		id       int64
		username string
		password string
	)

	if err := queryUser.Scan(&id, &username, &password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No existing user with that username",
		})
		return
	}

	if ok, _ := argon2id.ComparePasswordAndHash(user.Password, password); ok {
		accessToken, refreshToken, err := helpers.CreateToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "SOMETHING_BAD_HAPPENED",
				"message": "WHAT THE FUCK",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":          "OK",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Incorrect password",
		})
		return
	}
}
