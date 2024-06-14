package main

import (
	"fmt"
	"context"
	"os"
	"os/signal"	// provides a function to notify context, which is like a shared state for all concurrent processes


	"github.com/generalr3d/learn-go-api-chi-redis/application"	//	importing our own package in our own directory
)

func main(){
	app := application.New()	// instantiate

	contxt, cancel :=signal.NotifyContext(context.Background(),os.Interrupt)	//	takes in a context and signal, and returns another context if there is a system interrupt signal, in this case define a ROOT level context
	//	careful, should only be used for initialization and tests
	//	also returns a cancellation function for this ctx and all its child contexts
	defer cancel()

	err := app.Start(contxt)
	if err != nil {
		fmt.Println("Failed to start app:",err)
	}

	
}