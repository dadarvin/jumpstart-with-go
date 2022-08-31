package router

import (
	"entry_task/cmd/webservice/handler"
	m "entry_task/cmd/webservice/middleware"
	"github.com/julienschmidt/httprouter"
)

func userRouter(router *httprouter.Router) {
	// Initiate middleware
	//m := middleware.NewMiddleware(u)

	// GET
	router.GET("/get-profile/:id", m.IsAuthorized(handler.GetProfileFunc))
	router.GET("/get-profile-pict/:id", m.IsAuthorized(handler.GetProfilePictFunc))

	//POST
	router.POST("/register", handler.RegisterUserFunc)
	router.POST("/login", handler.LoginFunc)
	router.POST("/uploadprofilepict", m.IsAuthorized(handler.UploadProfilePictFunc))

	//PUT
	router.PUT("/change-nickname", m.IsAuthorized(handler.EditUserFunc))
}
