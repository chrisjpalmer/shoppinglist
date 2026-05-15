package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/chrisjpalmer/shoppinglist/backend/api"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

	spsrv, err := shopping.NewServer(planningSiteURL)
	if err != nil {
		log.Fatal("error during set up of shopping server", err)
	}

	srv := startServer(apisrv, spsrv, 8080)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	defer cancel()

	<-ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error while trying to close server: %v", err)
	}

	log.Println("server closed")
}

func startServer(apisrv *api.Server, spsrv *shopping.Server, port int) *http.Server {
	mux := http.NewServeMux()
	apisrv.RegisterRoutes(mux)
	spsrv.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error:", err)
		}
	}()

	log.Printf("server started, listening on :%d", port)

	return srv
}
