package server

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/db"
	"github.com/chrisjpalmer/shoppinglist/backend/gen"
	"github.com/chrisjpalmer/shoppinglist/backend/gen/genconnect"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	db *generated.Queries
}

func (s *Server) GetIngredients(ctx context.Context, rq *connect.Request[gen.GetIngredientsRequest]) (*connect.Response[gen.GetIngredientsResponse], error) {

}
func (s *Server) CreateIngredient(ctx context.Context, rq *connect.Request[gen.CreateIngredientRequest]) (*connect.Response[gen.CreateIngredientResponse], error) {

}
func (s *Server) UpdateIngredient(ctx context.Context, rq *connect.Request[gen.UpdateIngredientRequest]) (*connect.Response[gen.UpdateIngredientResponse], error) {

}
func (s *Server) DeleteIngredient(ctx context.Context, rq *connect.Request[gen.DeleteIngredientRequest]) (*connect.Response[gen.DeleteIngredientResponse], error) {

}

func New() (*Server, error) {
	db, err := db.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	return &Server{db: db}, nil
}

func (s *Server) Listen() error {
	mux := http.NewServeMux()
	path, handler := genconnect.NewShoppingListServiceHandler(s)
	mux.Handle(path, handler)

	return http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{}))
}
