package shopping

import (
	"context"
	"embed"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/page"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
	"golang.org/x/net/websocket"
)

//go:embed assets/*
var assets embed.FS

// Server - the server for giving you nice hello greetings
type Server struct {
	srv  http.Server
	done chan struct{}
}

// NewServer - creates a new server
func NewServer(port int) *Server {
	mux := http.NewServeMux()

	srv := &Server{
		srv: http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
		done: make(chan struct{}),
	}

	// serve one route on `/` which will be our hello page
	mux.Handle("/reload", &websocket.Server{
		Handler: srv.handleReload,
	})
	mux.HandleFunc("/", handleRootPage)
	mux.Handle("/assets/", http.FileServerFS(assets))
	mux.HandleFunc("/want", handleWantPage)
	mux.HandleFunc("/got", handleGotPage)
	mux.HandleFunc("/shop", handleShopPage)

	return srv
}

func (m *Server) handleReload(conn *websocket.Conn) {
	defer conn.Close()

	addr := conn.LocalAddr().String()

	fmt.Println("/reload: ", addr)

	done := make(chan struct{})

	go func() {
		buf := make([]byte, 100)
		defer close(done)

		for {
			_, err := conn.Read(buf)
			if err == io.EOF {
				return
			}
		}
	}()

	select {
	case <-done:
		fmt.Println("/reload: ", addr, "socket closed")
	case <-m.done:
		fmt.Println("/reload: ", addr, "sending closing signal")
	}

}

func handleRootPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/want")
	w.WriteHeader(http.StatusFound)
}

func handleWantPage(w http.ResponseWriter, r *http.Request) {
	pctx := page.NewContext(r)
	templ.Handler(render.WantPage(pctx)).ServeHTTP(w, r)
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
