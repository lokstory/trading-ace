package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"trading-ace/internal/api"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	app, err := api.NewApp(ctx, os.Getenv("CONFIG_PATH"))
	if err != nil {
		log.Panicln(err)
	}

	app.Start()

	defer func() {
		cancel()
		app.Close()
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	<-signalChan
}
