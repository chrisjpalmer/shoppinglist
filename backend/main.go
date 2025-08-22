package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/chrisjpalmer/shoppinglist/backend/server"
)

func main() {
	srv, err := server.New()
	if err != nil {
		log.Fatal("error during set up of server", err)
	}

	log.Println("server started, listening on 8080")

	go srv.Listen()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	defer cancel()

	<-ctx.Done()

	log.Println("server closed")
}
