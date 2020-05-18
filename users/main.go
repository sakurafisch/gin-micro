package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/web"
)

type User struct{}

var users = json.RawMessage(`[{"username": "winnerwinter", "email": "username@email.com"}, {"username": "bill", "email": "bill@biu.com"}]`)

func (user *User) List(context *gin.Context) {
	log.Print("Received User.list API request")
	context.JSON(http.StatusOK, gin.H{"result": true, "data": users})
}

func main() {
	user := new(User)

	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.NoMethod(func(context *gin.Context) {
		context.JSON(http.StatusMethodNotAllowed, gin.H{"result": false, "error": "Method Not Allowed"})
		return
	})
	router.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusNotFound, gin.H{"result": false, "error": "Endpoint Not Found"})
		return
	})
	router.GET("/user/list", user.List)

	service := web.NewService(
		web.Name("go.micro.api.user"),
	)
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	service.Handle("/", router)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
