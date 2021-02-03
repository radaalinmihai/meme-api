package profile

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"meme/db"
	"meme/helpers"
	"net/http"
	"strings"
)

type Profile struct {
	Id        int
	UserId    db.NullString `json:"userId "db:"userId"`
	Avatar    db.NullString `json:"avatar" db:"avatar"`
	FirstName db.NullString `json:"firstName" db:"firstName"`
	LastName  db.NullString `json:"lastName" db:"lastName"`
}

func GetProfile(c *gin.Context) {
	headerToken := c.Request.Header.Get("Authorization")
	token := strings.Split(headerToken, "Bearer ")
	tokenClaims, err := helpers.ParseToken(token[1])

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something bad happened",
		})
		return
	}

	sql := `
		SELECT avatar, firstName, lastName FROM profiles
		RIGHT JOIN users u ON profiles.userId=u.id
		WHERE username=?
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
		"code": "OK",
		"profile": gin.H{
			"avatar":    profile.Avatar.String,
			"firstName": profile.FirstName.String,
			"lastName":  profile.LastName.String,
		},
	})
}

func UpdateProfile(c *gin.Context) {
	profile := Profile{}
	profileId := c.Param("id")

	if err := c.ShouldBindJSON(&profile); err != nil {
		fmt.Println(err.Error())
		return
	}

	sql := "UPDATE profiles SET avatar=?, firstName=?, lastName=? WHERE id=?"

	_, err := db.MemeDB.Exec(sql, profile.Avatar.String, profile.FirstName.String, profile.LastName.String, profileId)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": "SOMETHING_HAPPENED",
			"message": err.Error(),
		})
		return
	}

	fmt.Println(profile.Id)

	c.JSON(http.StatusOK, gin.H{
		"code": "OK",
	})
}
