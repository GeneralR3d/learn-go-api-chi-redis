package main

import (
	"fmt"
	"net/http"

	"github.com/GeneralR3d/learn-go-api-chi-redis/application"	//	importing our own package in our own directory
)

func main(){
	app := application.New()	// instantiate

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Errorf("Failed to start app: %w",err)
	}
}