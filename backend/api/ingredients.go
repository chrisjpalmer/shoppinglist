package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/gen"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
)

func (s *Server) GetIngredients(ctx context.Context, rq *connect.Request[gen.GetIngredientsRequest]) (*connect.Response[gen.GetIngredientsResponse], error) {
	ii, err := s.sql.GetIngredients(ctx)
	if err != nil {
		return nil, err
	}

	var gii []*gen.Ingredient

	for _, i := range ii {
		gii = append(gii, &gen.Ingredient{
			Id:   i.ID,
			Name: i.Name,
		})
	}

	return connect.NewResponse(&gen.GetIngredientsResponse{Ingredients: gii}), nil
}
func (s *Server) CreateIngredient(ctx context.Context, rq *connect.Request[gen.CreateIngredientRequest]) (*connect.Response[gen.CreateIngredientResponse], error) {
	id, err := s.sql.CreateIngredient(ctx, rq.Msg.Ingredient.Name)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.CreateIngredientResponse{IngredientId: id}), nil
}
func (s *Server) UpdateIngredient(ctx context.Context, rq *connect.Request[gen.UpdateIngredientRequest]) (*connect.Response[gen.UpdateIngredientResponse], error) {
	err := s.sql.UpdateIngredient(ctx, generated.UpdateIngredientParams{
		ID:   rq.Msg.Ingredient.Id,
		Name: rq.Msg.Ingredient.Name,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.UpdateIngredientResponse{}), nil
}
func (s *Server) DeleteIngredient(ctx context.Context, rq *connect.Request[gen.DeleteIngredientRequest]) (*connect.Response[gen.DeleteIngredientResponse], error) {
	err := s.sql.DeleteIngredient(ctx, rq.Msg.IngredientId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.DeleteIngredientResponse{}), nil
}
