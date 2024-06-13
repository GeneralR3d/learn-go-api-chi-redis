package main

import (
	"fmt"
	"net/http"	// includes methods to create both http clients and servers

	"github.com/go-chi/chi/v5"	//	chi is a router which helps us route different URL paths
	"github.com/go-chi/chi/v5/middleware"
)

func main(){

	router := chi.NewRouter()	//	instatiate
	router.Use(middleware.Logger)	//	start logger
	router.Get("/hello",basicHandler)	//	this maps /hello path to our http handler function which is basicHandler
	

	server := &http.Server{		//	server is a pointer var to a http server
		Addr: ":3000",		//	server address
		//Handler: http.HandlerFunc(basicHandler),	//	need to define our own custom handler function
		Handler: router, //	the router itself conforms to 
	}

	err := server.ListenAndServe()	//start server
	if err != nil {
		fmt.Println("failed to listen to server",err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {		//	this is standard for all handler functions, a responseWriter w and a request pointer r
	w.Write([]byte("Hello world"))	//	cast string to byte array
}