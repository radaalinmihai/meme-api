package profile

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"meme/db"
	"meme/helpers"
	"net/http"
	"strings"
)

type Profile struct {
	Id        string
	Username  string          `json:"username" db:"username"`
	Email     string          `json:"email" db:"email"`
	CreatedAt string          `json:"created_at" db:"created_at"`
	UpdatedAt *sql.NullString `json:"updated_at" db:"updated_at"`
	ProfileId int             `json:"profileId" db:"profileId"`
	UserId    string          `json:"userId" db:"userId"`
	Avatar    string          `json:"avatar" db:"avatar"`
	FirstName string          `json:"firstName" db:"firstName"`
	LastName  string          `json:"lastName" db:"lastName"`
}

func Get(c *gin.Context) {
	headerToken := c.Request.Header.Get("Authorization")
	token := strings.Split(headerToken, "Bearer ")
	tokenClaims, err := helpers.ParseToken(token[1])
	fmt.Println(tokenClaims)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something bad happened",
		})
		return
	}

	sql := `
		SELECT
       	users.id,
       	users.username,
       	users.email,
       	users.created_at,
       	users.updated_at,
       	profiles.id as profileId,
       	profiles.userId,
       	profiles.avatar,
       	profiles.firstName,
       	profiles.lastName
		FROM users
		INNER JOIN profiles ON profiles.userId = users.id
		WHERE users.username=?;
	`
	profile := Profile{}
	if err := db.MemeDB.Get(&profile, sql, tokenClaims["username"]); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":  "NO_PROFILE",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "OK",
		"profile": profile,
	})
}

func Update(c *gin.Context) {
	profile := Profile{}
	profileId := c.Param("id")

	if err := c.ShouldBindJSON(&profile); err != nil {
		fmt.Println(err.Error())
		return
	}

	sql := "UPDATE profiles SET avatar=?, firstName=?, lastName=? WHERE id=?"

	_, err := db.MemeDB.Exec(sql, profile.Avatar, profile.FirstName, profile.LastName, profileId)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": err.Error(),
		})
		return
	}

	fmt.Println(profile.Id)

	c.JSON(http.StatusOK, gin.H{
		"code": "OK",
	})
}
