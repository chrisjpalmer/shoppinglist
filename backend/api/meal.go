package api

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/genpb"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
)

func (s *Server) GetMeals(ctx context.Context, rq *connect.Request[genpb.GetMealsRequest]) (*connect.Response[genpb.GetMealsResponse], error) {
	mm, err := s.sql.GetMeals(ctx)
	if err != nil {
		return nil, err
	}

	var gmm []*genpb.Meal

	for _, m := range mm {
		var ig []*genpb.IngredientRef
		err := unmarshalJSON(m.Ingredients, &ig)
		if err != nil {
			return nil, err
		}

		var previewImageUrl string
		s, ok := m.PreviewImageUrl.(string)
		if !ok {
			return nil, fmt.Errorf("error casting PreviewImageUrl to a string")
		}

		previewImageUrl = s

		var ingredientsImageUrl string
		s, ok = m.IngredientsImageUrl.(string)
		if !ok {
			return nil, fmt.Errorf("error casting IngredientsImageUrl to a string")
		}

		ingredientsImageUrl = s

		gmm = append(gmm, &genpb.Meal{
			Id:                  m.ID,
			Name:                m.Name,
			IngredientRefs:      ig,
			RecipeUrl:           m.RecipeUrl,
			PreviewImageUrl:     previewImageUrl,
			IngredientsImageUrl: ingredientsImageUrl,
		})
	}

	return connect.NewResponse(&genpb.GetMealsResponse{Meals: gmm}), nil
}
func (s *Server) CreateMeal(ctx context.Context, rq *connect.Request[genpb.CreateMealRequest]) (*connect.Response[genpb.CreateMealResponse], error) {
	igstr, err := marshalJSON(rq.Msg.Meal.IngredientRefs)
	if err != nil {
		return nil, err
	}

	id, err := s.sql.CreateMeal(ctx, gensql.CreateMealParams{
		Name:                rq.Msg.Meal.Name,
		Ingredients:         igstr,
		RecipeUrl:           rq.Msg.Meal.RecipeUrl,
		PreviewImageUrl:     rq.Msg.Meal.PreviewImageUrl,
		IngredientsImageUrl: rq.Msg.Meal.IngredientsImageUrl,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.CreateMealResponse{MealId: id}), nil
}
func (s *Server) UpdateMeal(ctx context.Context, rq *connect.Request[genpb.UpdateMealRequest]) (*connect.Response[genpb.UpdateMealResponse], error) {
	igstr, err := marshalJSON(rq.Msg.Meal.IngredientRefs)
	if err != nil {
		return nil, err
	}

	err = s.sql.UpdateMeal(ctx, gensql.UpdateMealParams{
		ID:                  rq.Msg.Meal.Id,
		Name:                rq.Msg.Meal.Name,
		Ingredients:         igstr,
		RecipeUrl:           rq.Msg.Meal.RecipeUrl,
		PreviewImageUrl:     rq.Msg.Meal.PreviewImageUrl,
		IngredientsImageUrl: rq.Msg.Meal.IngredientsImageUrl,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.UpdateMealResponse{}), nil
}
func (s *Server) DeleteMeal(ctx context.Context, rq *connect.Request[genpb.DeleteMealRequest]) (*connect.Response[genpb.DeleteMealResponse], error) {
	err := s.sql.DeleteMeal(ctx, rq.Msg.MealId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.DeleteMealResponse{}), nil
}
