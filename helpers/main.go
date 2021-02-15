package helpers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateToken(username string) (string, string, error) {
	tokenSecret := os.Getenv("TOKEN_SECRET")

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"expire": time.Now().Local().Add(time.Hour * 3),
	})

	accessTokenString, err := accessToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_token": accessTokenString,
		"expire": time.Now().Local().Add(time.Hour * 24 * 90),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	tokenSecret := os.Getenv("TOKEN_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func RefreshToken(c *gin.Context) {
	bearer := c.Request.Header.Get("Authorization")
	refreshToken := strings.Split(bearer, " ")

	if len(refreshToken) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "MISSING_TOKEN",
			"message": "Missing refresh token",
		})
		return
	}

	refreshTokenClaims, err := ParseToken(refreshToken[1])
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "REFRESH_TOKEN_EXPIRED",
			"message": "Ceva nu a mers bine, te rugam sa te loghezi din nou",
		})
		return
	}

	refreshTokenExpire, _ := time.Parse(time.RFC3339, refreshTokenClaims["expire"].(string))
	elapsedTime := refreshTokenExpire.Unix() - time.Now().Unix()

	if elapsedTime < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "REFRESH_TOKEN_EXPIRED",
			"message": "Ceva nu a mers bine, te rugam sa te loghezi din nou",
		})
		return
	}

	accessTokenClaims, err := ParseToken(refreshTokenClaims["access_token"].(string))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "REFRESH_TOKEN_EXPIRED",
			"message": "Ceva nu a mers bine, te rugam sa te loghezi din nou",
		})
		return
	}

	newAccessToken, newRefreshToken, err := CreateToken(accessTokenClaims["username"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "REFRESH_TOKEN_EXPIRED",
			"message": "Ceva nu a mers bine, te rugam sa te loghezi din nou",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "OK",
		"access_token": newAccessToken,
		"refresh_token": newRefreshToken,
	})
}