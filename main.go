package main

import (
	"fmt"
	auth "jwt-todo/auth/jwt"
	"log"
	"net/http"
	"strings"

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

type Todo struct {
	UserID uint64 `json:"user_id"`
	Title  string `json:"title"`
}

type Header struct {
	Authorization string `header:"Authorization"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ExtractToken(token string) string {
	t := strings.Split(token, " ")
	if len(t) == 2 {
		return t[1]
	}

	return ""
}

func main() {
	r := gin.Default()

	gin.SetMode(gin.DebugMode)
	r.POST("/api/login", Login)
	r.GET("/api/test", func(c *gin.Context) {
		var h Header
		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(http.StatusUnauthorized, "Invalid Authorization")
			return
		}

		token := ExtractToken(h.Authorization)
		c.JSON(http.StatusOK, token)
	})

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
	err = jwt.SetTokenRedis(user.ID, &token)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	c.JSON(http.StatusOK, token)
}
