package router

import (
	"github.com/julienschmidt/httprouter"
)

func Init() *httprouter.Router {
	// Initialize a new router
	router := httprouter.New()

	// Define each handler here
	userRouter(router)

	return router
}
