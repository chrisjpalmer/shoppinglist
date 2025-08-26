package server

import (
	"context"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/gen"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
)

func (s *Server) GetMeals(ctx context.Context, rq *connect.Request[gen.GetMealsRequest]) (*connect.Response[gen.GetMealsResponse], error) {
	mm, err := s.sql.GetMeals(ctx)
	if err != nil {
		return nil, err
	}

	var gmm []*gen.Meal

	for _, m := range mm {
		var ig []*gen.IngredientRef
		err := unmarshalJSON(m.Ingredients, &ig)
		if err != nil {
			return nil, err
		}

		gmm = append(gmm, &gen.Meal{
			Id:             m.ID,
			Name:           m.Name,
			IngredientRefs: ig,
			RecipeUrl:      m.RecipeUrl,
		})
	}

	return connect.NewResponse(&gen.GetMealsResponse{Meals: gmm}), nil
}
func (s *Server) CreateMeal(ctx context.Context, rq *connect.Request[gen.CreateMealRequest]) (*connect.Response[gen.CreateMealResponse], error) {
	igstr, err := marshalJSON(rq.Msg.Meal.IngredientRefs)
	if err != nil {
		return nil, err
	}

	id, err := s.sql.CreateMeal(ctx, generated.CreateMealParams{
		Name:        rq.Msg.Meal.Name,
		Ingredients: igstr,
		RecipeUrl:   rq.Msg.Meal.RecipeUrl,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.CreateMealResponse{MealId: id}), nil
}
func (s *Server) UpdateMeal(ctx context.Context, rq *connect.Request[gen.UpdateMealRequest]) (*connect.Response[gen.UpdateMealResponse], error) {
	igstr, err := marshalJSON(rq.Msg.Meal.IngredientRefs)
	if err != nil {
		return nil, err
	}

	err = s.sql.UpdateMeal(ctx, generated.UpdateMealParams{
		ID:          rq.Msg.Meal.Id,
		Name:        rq.Msg.Meal.Name,
		Ingredients: igstr,
		RecipeUrl:   rq.Msg.Meal.RecipeUrl,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.UpdateMealResponse{}), nil
}
func (s *Server) DeleteMeal(ctx context.Context, rq *connect.Request[gen.DeleteMealRequest]) (*connect.Response[gen.DeleteMealResponse], error) {
	err := s.sql.DeleteMeal(ctx, rq.Msg.MealId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.DeleteMealResponse{}), nil
}
