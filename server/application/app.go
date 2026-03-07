package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"github.com/klados/weather_monitor/internal/config"
)

type App struct {
	router http.Handler
	fireDb *db.Client
	server *http.Server
	config *config.Config
}

func New() *App {
	router := loadRoutes()

	return &App{
		router: router,
	}
}

func (a *App) Start(ctx context.Context) error {

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load env variables: %w", err)
	}
	
	conf := &firebase.Config{
		ProjectID:   cfg.Firebase.ProjectID,
		DatabaseURL: cfg.Firebase.DatabaseURL,
	}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return fmt.Errorf("failed to initialize firebase: %w", err)
	}

	fireDb, err := app.Database(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	a.config = cfg
	
	a.fireDb = fireDb

	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%s",cfg.Port),
		Handler:      a.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("server running on", cfg.Port)

	ch := make(chan error, 1)
	
	go func() {	
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
	
		close(ch)
	}()
	

	select {
		case err = <-ch:
			return err
		case <-ctx.Done():
		 	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return a.server.Shutdown(shutdownCtx)
	}
	
}
