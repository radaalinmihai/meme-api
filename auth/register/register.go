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

	tx, err := db.MemeDB.Begin()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "INTERNAL_ERROR",
			"errors": err.Error(),
		})
		return
	}

	{
		tx.QueryRow("SET @uuid=uuid()")
		sql := `
			INSERT INTO users(id, username, password, email)
			VALUES(@uuid, ?, ?, ?)
		`
		if _, err := tx.Exec(sql, newUser.Username, hashedPassword, newUser.Email); err != nil {
			_ = tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": "INTERNAL_ERROR",
				"errors": err.Error(),
			})
			return
		}

	}

	{
		sql := `
			INSERT INTO profiles(userId)
			VALUES(@uuid)
		`
		if _, err := tx.Exec(sql); err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code": "INTERNAL_ERROR",
				"errors": err.Error(),
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		fmt.Print(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "INTERNAL_ERROR",
			"errors": err.Error(),
		})
		return
	}


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
