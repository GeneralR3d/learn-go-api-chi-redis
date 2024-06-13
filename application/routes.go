package application

import (
	"context"
	"net/http"	// includes methods to create both http clients and servers

	"github.com/go-chi/chi/v5"	//	chi is a router which helps us route different URL paths
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()	//	instatiate

	router.Use(middleware.Logger)	//	start logger

	router.Get("/",func (w http.ResponseWriter, r *http.Request){	// default / route, replaces the basic handler thing w an anon function
		w.WriteHeader(http.StatusOk)
	})

	return router
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