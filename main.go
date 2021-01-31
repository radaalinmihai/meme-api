package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"meme/api/profile"
	memeAuth "meme/auth"
	"meme/auth/login"
	"meme/auth/register"
	"meme/db"
	"meme/helpers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file :(")
	}
	db.ConnectDB()
	r := gin.New()

	api := r.Group("/api")

	auth := api.Group("/auth")
	auth.POST("/login", login.Handler)
	auth.POST("/register", register.Handler)
	auth.GET("/refreshToken", helpers.RefreshToken)

	api.Use(memeAuth.Middleware())
	api.GET("/profile", profile.GetProfile)
	api.PUT("/profile/:id", profile.UpdateProfile)

	_ = r.Run(":4000")
}