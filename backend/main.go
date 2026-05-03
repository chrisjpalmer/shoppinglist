package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chrisjpalmer/shoppinglist/backend/api"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping"
)

func main() {
	planningSiteURL := os.Getenv("PLANNING_SITE_URL")
	if planningSiteURL == "" {
		planningSiteURL = "http://localhost:3000"
	}

	apisrv, err := api.NewServer()
	if err != nil {
		log.Fatal("error during set up of api server", err)
	}

	spsrv, err := shopping.NewServer(8081, planningSiteURL)
	if err != nil {
		log.Fatal("error during set up of shopping server", err)
	}

	go apisrv.Listen()

	log.Println("api server started, listening on 8080")

	go spsrv.Listen()

	log.Println("shopping server started, listening on 8081")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	defer cancel()

	<-ctx.Done()

	if err := spsrv.Close(); err != nil {
		log.Printf("error while trying to close shopping server: %v", err)
	}

	log.Println("server closed")
}
