package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	gofirestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/klados/weather_monitor/internal/config"
	"google.golang.org/api/option"
)

type App struct {
	router    http.Handler
	fireStore *gofirestore.Client
	server    *http.Server
	config    *config.Config
}

func New() *App {
	return &App{}
}

func (a *App) Start(ctx context.Context) error {

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load env variables: %w", err)
	}

	opt := option.WithCredentialsFile(cfg.Firebase.CredentialsFile)
	fbConfig := &firebase.Config{
		ProjectID: cfg.Firebase.ProjectID,
	}

	app, err := firebase.NewApp(context.Background(), fbConfig, opt)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	fireStore, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize firestore: %w", err)
	}

	a.config = cfg
	a.fireStore = fireStore
	a.router = loadRoutes(a.fireStore)

	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
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
