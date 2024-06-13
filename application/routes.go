package application

import (
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

