package server

import (
	"context"
	"net/http"

	"github.com/chrisjpalmer/shoppinglist/backend/gen/genconnect"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
	"github.com/chrisjpalmer/shoppinglist/backend/sql"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	sql *generated.Queries
}

func New() (*Server, error) {
	sql, err := sql.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	return &Server{sql: sql}, nil
}

func (s *Server) Listen() error {
	mux := http.NewServeMux()
	path, handler := genconnect.NewShoppingListServiceHandler(s)
	mux.Handle(path, handler)

	return http.ListenAndServe(":8080", h2c.NewHandler(cors.Default().Handler(mux), &http2.Server{}))
}
