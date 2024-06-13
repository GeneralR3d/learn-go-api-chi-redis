package application

import(
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