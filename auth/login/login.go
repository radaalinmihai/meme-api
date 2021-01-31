package login

import (
	"database/sql"
	"fmt"
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

type User struct {
	Id        string
	Username  string
	Password  string
	Email     string
	CreatedAt sql.NullString `db:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at"`
}

func Handler(c *gin.Context) {
	user := Body{}
	dbUser := User{}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sql := "SELECT * FROM users WHERE username=?"
	queryUser := db.MemeDB.QueryRowx(sql, user.Name)

	if err := queryUser.StructScan(&dbUser); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No existing user with that username",
		})
		return
	}

	if ok, _ := argon2id.ComparePasswordAndHash(user.Password, dbUser.Password); ok {
		accessToken, refreshToken, err := helpers.CreateToken(dbUser.Username)
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
