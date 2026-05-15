package shopping

import (
	"context"
	"embed"
	"net/http"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/page"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
	"github.com/chrisjpalmer/shoppinglist/backend/sql"
)

const (
	fragmentParam = "fragment"
)

//go:embed assets/*
var assets embed.FS

type Server struct {
	planningSiteURL string
	sql             *gensql.Queries
}

func NewServer(planningSiteURL string) (*Server, error) {
	sql, err := sql.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	return &Server{
		planningSiteURL: planningSiteURL,
		sql:             sql,
	}, nil
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/shopping", handleRootPage)
	mux.Handle("/shopping/assets/", http.StripPrefix("/shopping", http.FileServerFS(assets)))
	mux.HandleFunc("/shopping/want", s.handleWantPage)
	mux.HandleFunc("/shopping/got", s.handleGotPage)
	mux.HandleFunc("/shopping/shop", s.handleShopPage)
}

func handleRootPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/shopping/want")
	w.WriteHeader(http.StatusFound)
}

func (s *Server) handleGotPage(w http.ResponseWriter, r *http.Request) {
	pctx := s.pageContext(r)
	templ.Handler(render.GotPage(pctx)).ServeHTTP(w, r)
}

func (s *Server) handleShopPage(w http.ResponseWriter, r *http.Request) {
	pctx := s.pageContext(r)
	templ.Handler(render.ShopPage(pctx)).ServeHTTP(w, r)
}

func (s *Server) pageContext(r *http.Request) page.Context {
	return page.NewContext(r, s.planningSiteURL)
}
