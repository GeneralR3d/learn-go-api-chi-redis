package application

import (
	"net/http"	// includes methods to create both http clients and servers
	"github.com/generalr3d/learn-go-api-chi-redis/handlers"
	"github.com/go-chi/chi/v5"	//	chi is a router which helps us route different URL paths
	"github.com/go-chi/chi/v5/middleware"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()	//	instatiate

	router.Use(middleware.Logger)	//	start logger

	router.Get("/",func (w http.ResponseWriter, r *http.Request){	// default / route, replaces the basic handler thing w an anon function
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders",loadOrderRoutes)	//	route to sub-router, dont service immediately


	return router
}

func loadOrderRoutes(router chi.Router){	//	sub router, will only receive all the /order routes

	orderHandler := &handlers.Order{}	//	instiate order

	router.Post("/",orderHandler.Create)
	router.Get("/",orderHandler.List)
	router.Get("/{id}",orderHandler.GetByID)	//	id field to extract the path parameter
	router.Put("/{id}",orderHandler.UpdateByID)
	router.Delete("/{id}",orderHandler.DeleteByID)
}

