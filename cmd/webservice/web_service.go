package webservice

import (
	"context"
	"entry_task/cmd/webservice/handler"
	"entry_task/cmd/webservice/router"
	"entry_task/internal/config"
	"entry_task/internal/repo"
	"entry_task/internal/usecase"
	"log"
	"net/http"
	"time"
)

//var (
//	srv *http.Server
//)

type Server struct {
	srv     *http.Server
	handler *handler.Handler
}

func Start() (stopFunc func()) {
	conf := config.Get()

	// Initialize Repository
	userRepo := repo.InitDependencies()

	// Initialize Usecase
	userUC := usecase.InitDependencies(userRepo)

	// Initialize Router for Handler
	router := router.Init()

	httpSrv := &http.Server{
		Addr:    ":" + conf.HttpPort,
		Handler: router,
	}

	server := &Server{
		srv:     httpSrv,
		handler: handler.New(userUC),
	}

	go func() {
		err := server.srv.ListenAndServe()
		if err != nil {
			log.Fatalln("Failed starting server on ", server.srv.Addr, err)
		}
	}()
	return func() {
		GracefulStop(server)
	}
}

func GracefulStop(server *Server) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), (15 * time.Second))
	defer cancel()

	log.Println("shuting down web on", server.srv.Addr)
	err = server.srv.Shutdown(ctx)
	if err != nil {
		log.Fatalln("failed shutdown server", err)
	}
	log.Println("web gracefully stopped")
	return
}
