package register

import (
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"log"
	"meme/db"
	"meme/helpers"
	"net/http"
)

type Body struct {
	Username string `form:"username" json:"username" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Handler(c *gin.Context) {
	newUser := Body{}

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "MISSING_FIELDS",
			"message": "Something is missing here",
		})
		return
	}

	hashedPassword, err := argon2id.CreateHash(newUser.Password, &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  1,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "FAILED_HASHING",
		})
		return
	}

	sql := "INSERT INTO users(username, email, password) VALUES(?, ?, ?)"
	res, err := db.MemeDB.Exec(sql, newUser.Username, newUser.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "FAILED_INSERTION",
		})
		panic(err.Error())
		return
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Println(rowsAffected)

	accessToken, refreshToken, err := helpers.CreateToken(newUser.Username)
	if err != nil {
		log.Fatal(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          "OK",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
