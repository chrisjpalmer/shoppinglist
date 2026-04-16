package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/genpb"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
)

func (s *Server) GetIngredients(ctx context.Context, rq *connect.Request[genpb.GetIngredientsRequest]) (*connect.Response[genpb.GetIngredientsResponse], error) {
	ii, err := s.sql.GetIngredients(ctx)
	if err != nil {
		return nil, err
	}

	var gii []*genpb.Ingredient

	for _, i := range ii {
		gii = append(gii, &genpb.Ingredient{
			Id:   i.ID,
			Name: i.Name,
		})
	}

	return connect.NewResponse(&genpb.GetIngredientsResponse{Ingredients: gii}), nil
}
func (s *Server) CreateIngredient(ctx context.Context, rq *connect.Request[genpb.CreateIngredientRequest]) (*connect.Response[genpb.CreateIngredientResponse], error) {
	id, err := s.sql.CreateIngredient(ctx, rq.Msg.Ingredient.Name)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.CreateIngredientResponse{IngredientId: id}), nil
}
func (s *Server) UpdateIngredient(ctx context.Context, rq *connect.Request[genpb.UpdateIngredientRequest]) (*connect.Response[genpb.UpdateIngredientResponse], error) {
	err := s.sql.UpdateIngredient(ctx, gensql.UpdateIngredientParams{
		ID:   rq.Msg.Ingredient.Id,
		Name: rq.Msg.Ingredient.Name,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.UpdateIngredientResponse{}), nil
}
func (s *Server) DeleteIngredient(ctx context.Context, rq *connect.Request[genpb.DeleteIngredientRequest]) (*connect.Response[genpb.DeleteIngredientResponse], error) {
	err := s.sql.DeleteIngredient(ctx, rq.Msg.IngredientId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.DeleteIngredientResponse{}), nil
}
