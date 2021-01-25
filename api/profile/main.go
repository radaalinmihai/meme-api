package profile

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"meme/db"
	"meme/helpers"
	"net/http"
	"strings"
)

func Handler(c *gin.Context) {
	headerToken := c.Request.Header.Get("Authorization")
	token := strings.Split(headerToken, "Bearer ")
	tokenClaims, err := helpers.ParseToken(token[1])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something bad happened",
		})
		return
	}

	sql := "SELECT * from users WHERE username=?"
	queryUser := db.MemeDB.QueryRow(sql, tokenClaims["username"])

	var (
		username string
		email    string
	)

	_ = queryUser.Scan(&username, &email)
	fmt.Print(queryUser.Err())

	fmt.Printf("Email: %s", email)
	c.JSON(http.StatusOK, gin.H{
		"code": "OK",
		"profile": gin.H{
			"username": username,
			"email":    email,
		},
	})
}
