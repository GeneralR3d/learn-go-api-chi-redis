package application

import(
	"fmt"
	"context"
	"net/http"

)

type App struct{
	router http.Handler
}

//	create constructor
func New() *App{
	app := &App{
		router:loadRoutes(),
	}

	return app	
}

func (a *App) Start(context context.Context) error { 	//	receiver, or method that belongs to App struct


	server := &http.Server{	//	server is a pointer var to a http server
		Addr: ":3000",
		Handler: a.router,	//	this.router for the App object
	}

	err := server.ListenAndServe()	//start server
	if err != nil{
		return fmt.Errorf("Failed to start server: %w ", err)	//	 error wrapping, wrap an error w another error around it and return that wrapped error
	}

	return nil	//	when there is no error, success
}