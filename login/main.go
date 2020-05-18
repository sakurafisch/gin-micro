package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/web"
)

type User struct{}

func (user *User) Login(context *gin.Context) {
	log.Print("Received Login API request")
	var request map[string]interface{}

	if err := json.NewDecoder(context.Request.Body).Decode(&request); err != nil {
		context.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	if request["username"] == nil {
		err := errors.New("Field username is required.")
		context.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	if request["password"] == nil {
		err := errors.New("Field password is required.")
		context.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"result": false, "error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"result": true, "data": request})
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
	router.POST("/login", user.Login)

	service := web.NewService(
		web.Name("go.Micro.api.login"),
	)
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	service.Handle("/", router)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
