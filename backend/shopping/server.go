package shopping

import (
	"context"
	"embed"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/page"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
	"github.com/chrisjpalmer/shoppinglist/backend/sql"
)

//go:embed assets/*
var assets embed.FS

// Server - the server for giving you nice hello greetings
type Server struct {
	srv  http.Server
	done chan struct{}
	sql  *generated.Queries
}

// NewServer - creates a new server
func NewServer(port int) (*Server, error) {
	mux := http.NewServeMux()

	sql, err := sql.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	srv := &Server{
		srv: http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
		done: make(chan struct{}),
		sql:  sql,
	}

	// serve one route on `/` which will be our hello page
	mux.HandleFunc("/", handleRootPage)
	mux.Handle("/assets/", http.FileServerFS(assets))
	mux.HandleFunc("/want", srv.handleWantPage)
	mux.HandleFunc("/got", handleGotPage)
	mux.HandleFunc("/shop", handleShopPage)

	return srv, nil
}

func handleRootPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/want")
	w.WriteHeader(http.StatusFound)
}

func handleGotPage(w http.ResponseWriter, r *http.Request) {
	pctx := page.NewContext(r)
	templ.Handler(render.GotPage(pctx)).ServeHTTP(w, r)
}

func handleShopPage(w http.ResponseWriter, r *http.Request) {
	pctx := page.NewContext(r)
	templ.Handler(render.ShopPage(pctx)).ServeHTTP(w, r)
}

// Listen - starts the server
func (s *Server) Listen() error {
	return s.srv.ListenAndServe()
}

// Close - gracefully closes the server
func (s *Server) Close() error {
	close(s.done)

	return s.srv.Shutdown(context.Background())
}
