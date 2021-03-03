package main

import (
	"fmt"
	auth "jwt-todo/auth/jwt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var user = User{
	ID:       1,
	Name:     "username",
	Password: "password",
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	r := gin.Default()

	gin.SetMode(gin.DebugMode)
	r.POST("/api/login", Login)

	log.Fatal(r.Run(":1234"))
}

func Login(c *gin.Context) {
	jwt, err := auth.NewJWTHandler()
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if user.Name != u.Name || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
	}

	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}
