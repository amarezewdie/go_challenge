package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	person_service "github.com/izymalhaw/go-crud/yishakterefe/internal/services/person"
)

// Server represents the HTTP server for the application.
type Server struct {
	port          int
	router        *http.ServeMux
	logger        *slog.Logger
	PersonService person_service.PersonServiceAbstrcatImpl
}

// NewApp initializes a new Server instance with the provided port, person service, and logger.
func NewApp(port int, personService person_service.PersonServiceAbstrcatImpl, logger *slog.Logger) *Server {

	app := &Server{
		router:        http.NewServeMux(),
		port:          port,
		PersonService: personService,
		logger:        logger,
	}

	app.Routes()
	wrappedRouter := app.enableCORS(app.router)
	app.router = http.NewServeMux()
	app.router.Handle("/", wrappedRouter)
	return app
}


func (app *Server) Run() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", strconv.Itoa(app.port)),
		Handler:      app.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit

		app.logger.Info("shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()
	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.Info("server stopped")
	return nil
}
