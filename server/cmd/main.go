package main

import (
	"context"
	"fmt"

	"github.com/klados/weather_monitor/application"
)

func main() {
	fmt.Println("start of magic")

	app := application.New()

	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("failed to start app", err)
	}
}
