package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"trading-ace/internal/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app, err := worker.NewApp(ctx, os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Panicln(err)
	}

	appErr := app.Start()
	if appErr != nil {
		log.Panicln(appErr)
	}

	defer func() {
		cancel()
		app.Close()
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	<-signalChan
}
