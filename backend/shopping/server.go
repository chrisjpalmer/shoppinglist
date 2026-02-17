package shopping

import (
	"context"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
)

// Server - the server for giving you nice hello greetings
type Server struct {
	srv http.Server
}

// NewServer - creates a new server
func NewServer(port int) *Server {
	mux := http.NewServeMux()

	// serve one route on `/` which will be our hello page
	mux.HandleFunc("/", handleHomePage)

	return &Server{
		srv: http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
	}
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	var (
		name = ""
		opts []func(*templ.ComponentHandler)
	)

	templ.Handler(render.HomePage(name), opts...).ServeHTTP(w, r)
}

// Listen - starts the server
func (s *Server) Listen() error {
	return s.srv.ListenAndServe()
}

// Close - gracefully closes the server
func (s *Server) Close() error {
	return s.srv.Shutdown(context.Background())
}
