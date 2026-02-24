package shopping

import (
	"context"
	"embed"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
)

//go:embed assets/*
var assets embed.FS

// Server - the server for giving you nice hello greetings
type Server struct {
	srv http.Server
}

// NewServer - creates a new server
func NewServer(port int) *Server {
	mux := http.NewServeMux()

	// serve one route on `/` which will be our hello page
	mux.HandleFunc("/", handleRootPage)
	mux.Handle("/assets/", http.FileServerFS(assets))
	mux.HandleFunc("/want", handleWantPage)
	mux.HandleFunc("/got", handleGotPage)
	mux.HandleFunc("/shop", handleShopPage)

	return &Server{
		srv: http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
	}
}

func handleRootPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/want")
	w.WriteHeader(http.StatusFound)
}

func handleWantPage(w http.ResponseWriter, r *http.Request) {
	templ.Handler(render.WantPage()).ServeHTTP(w, r)
}

func handleGotPage(w http.ResponseWriter, r *http.Request) {
	templ.Handler(render.GotPage()).ServeHTTP(w, r)
}

func handleShopPage(w http.ResponseWriter, r *http.Request) {
	templ.Handler(render.ShopPage()).ServeHTTP(w, r)
}

// Listen - starts the server
func (s *Server) Listen() error {
	return s.srv.ListenAndServe()
}

// Close - gracefully closes the server
func (s *Server) Close() error {
	return s.srv.Shutdown(context.Background())
}
