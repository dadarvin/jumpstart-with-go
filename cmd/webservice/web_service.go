package webservice

import (
	"context"
	"entry_task/cmd/webservice/handler"
	"entry_task/cmd/webservice/middleware"
	"entry_task/internal/config"
	"entry_task/internal/repo"
	"entry_task/internal/usecase"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

//var (
//	srv *http.Server
//)

type Server struct {
	srv        *http.Server
	handler    *handler.Handler
	middleware *middleware.Middleware
}

func Init() (stopFunc func()) {
	conf := config.Get()

	// Initialize Repository
	userRepo := repo.InitDependencies()

	// Initialize Usecase
	userUC := usecase.InitDependencies(userRepo)

	// Initialize Router for Handler
	router := httprouter.New()

	httpSrv := &http.Server{
		Addr:    ":" + conf.HttpPort,
		Handler: router,
	}

	server := &Server{
		srv:        httpSrv,
		handler:    handler.New(userUC),
		middleware: middleware.New(conf.AuthConfig.JWTSecret),
	}
	server.Start(router)

	go func() {
		log.Println("Server running on", conf.HttpPort)
		err := server.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("failed starting web on", server.srv.Addr, err)
		}
	}()
	return func() {
		server.GracefulStop()
	}
}

func (s *Server) Start(router *httprouter.Router) {
	// GET
	router.GET("/get-profile/:id", s.middleware.IsAuthorized(s.handler.GetProfileFunc()))
	router.GET("/get-profile-pict/:id", s.middleware.IsAuthorized(s.handler.GetProfilePictFunc()))

	//POST
	router.POST("/register", s.handler.RegisterUserFunc())
	router.POST("/login", s.handler.LoginFunc())
	router.POST("/uploadprofilepict", s.middleware.IsAuthorized(s.handler.UploadProfilePictFunc()))

	//PUT
	router.PUT("/change-nickname", s.middleware.IsAuthorized(s.handler.EditUserFunc()))
}

func (s *Server) GracefulStop() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), (15 * time.Second))
	defer cancel()

	log.Println("shuting down web on", s.srv.Addr)
	err = s.srv.Shutdown(ctx)
	if err != nil {
		log.Fatalln("failed shutdown server", err)
	}
	log.Println("web gracefully stopped")
	return
}
