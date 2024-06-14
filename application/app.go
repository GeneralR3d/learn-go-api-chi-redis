package application

import(
	"fmt"
	"context"	//	provides mechanism for communicating and handling cancellations
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct{
	router http.Handler
	rdb *redis.Client
}

//	create constructor
func New() *App{
	app := &App{
		router:loadRoutes(),
		rdb: redis.NewClient(&redis.Options{}),
	}

	return app	
}

func (a *App) Start(contxt context.Context) error { 	//	receiver, or method that belongs to App struct


	server := &http.Server{	//	server is a pointer var to a http server
		Addr: ":3000",
		Handler: a.router,	//	this.router for the App object
	}

	//	wrap the closing of redis connection in a anon function, as by itself it doesnt work in a function that returns an error
	defer func() {
		if err:=a.rdb.Close() ; err != nil {
			fmt.Errorf("failed to close redis: %w",err)
		}
	}()

	err := a.rdb.Ping(contxt).Err()
	if err != nil{
		return fmt.Errorf("failed to connect to redis: %w",err)
	}

	fmt.Println("starting server...")

	channel := make(chan error, 1)	// uffered channel (non-blocking) of size 1 containing an error

	go func(){	// wrap the server start in a anon go-routine, which ensures it does not block main thread
		err = server.ListenAndServe()	//start server
		if err != nil{
			channel <- fmt.Errorf("failed to start server: %w ", err)	//	 error wrapping, wrap an error w another error around it and return that wrapped error
		}
	close(channel)	//	informs any listener to stop expecting any more values from the channel
	}()

	// err, open := <- channel	//	listen from the channel and determine if its closed or not
	// if !open {
	// 	//do something
	// }

	select {	//	select statement is like switch but for channels, allows us to block on (listen to) multiple channels and the first one with a value will proceed and only that one will execute
	case err = <-channel: 
		return err
	case <- contxt.Done():		//	also returns a channel
		// return server.ShutDown(context.Background())	//	shutdown could run indefinitely
		timeout, cancel := context.WithTimeout(context.Background(),time.Second *10)
		defer cancel()
		return server.Shutdown(timeout)	//http.Server.Shutdown()

	}



	return nil	//	when there is no error, success
}