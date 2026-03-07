package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/klados/weather_monitor/application"
)

func main() {
	fmt.Println("start of magic")

	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app", err)
	}
}
