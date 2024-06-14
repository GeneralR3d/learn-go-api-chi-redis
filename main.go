package main

import (
	"fmt"
	"context"
	"github.com/generalr3d/learn-go-api-chi-redis/application"	//	importing our own package in our own directory
)

func main(){
	app := application.New()	// instantiate

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("Failed to start app:",err)
	}
}