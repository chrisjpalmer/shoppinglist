package api

import (
	"context"
	"net/http"

	"github.com/chrisjpalmer/shoppinglist/backend/genpb/genpbconnect"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
	"github.com/chrisjpalmer/shoppinglist/backend/sql"
	"github.com/rs/cors"
)

type Server struct {
	sql *gensql.Queries
}

func NewServer() (*Server, error) {
	sql, err := sql.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	return &Server{sql: sql}, nil
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	apiMux := http.NewServeMux()

	path, handler := genpbconnect.NewShoppingListServiceHandler(s)

	apiMux.Handle(path, handler)

	apiMux.HandleFunc("/meal/{id}/image/preview", s.handleMealPreviewImage)
	apiMux.HandleFunc("/meal/{id}/image/ingredients", s.handleMealIngredientsImage)

	mux.Handle("/api/", cors.AllowAll().Handler(http.StripPrefix("/api", apiMux)))
}
