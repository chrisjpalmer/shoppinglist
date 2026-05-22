package api

import (
	"context"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/genpb"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
)

func (s *Server) GetIngredientCategories(ctx context.Context, rq *connect.Request[genpb.GetIngredientCategoriesRequest]) (*connect.Response[genpb.GetIngredientCategoriesResponse], error) {
	cc, err := s.sql.GetIngredientCategories(ctx)
	if err != nil {
		return nil, err
	}

	var gcc []*genpb.IngredientCategory

	for _, c := range cc {
		gcc = append(gcc, &genpb.IngredientCategory{
			Id:   c.ID,
			Name: c.Name,
		})
	}

	return connect.NewResponse(&genpb.GetIngredientCategoriesResponse{IngredientCategories: gcc}), nil
}

func (s *Server) CreateIngredientCategory(ctx context.Context, rq *connect.Request[genpb.CreateIngredientCategoryRequest]) (*connect.Response[genpb.CreateIngredientCategoryResponse], error) {
	cat, err := s.sql.CreateIngredientCategory(ctx, rq.Msg.IngredientCategory.Name)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.CreateIngredientCategoryResponse{IngredientCategoryId: cat.ID, SortIndex: cat.SortIndex}), nil
}

func (s *Server) UpdateIngredientCategory(ctx context.Context, rq *connect.Request[genpb.UpdateIngredientCategoryRequest]) (*connect.Response[genpb.UpdateIngredientCategoryResponse], error) {
	err := s.sql.UpdateIngredientCategory(ctx, gensql.UpdateIngredientCategoryParams{
		ID:   rq.Msg.IngredientCategory.Id,
		Name: rq.Msg.IngredientCategory.Name,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.UpdateIngredientCategoryResponse{}), nil
}

func (s *Server) DeleteIngredientCategory(ctx context.Context, rq *connect.Request[genpb.DeleteIngredientCategoryRequest]) (*connect.Response[genpb.DeleteIngredientCategoryResponse], error) {
	err := s.sql.DeleteIngredientCategory(ctx, rq.Msg.IngredientCategoryId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.DeleteIngredientCategoryResponse{}), nil
}

func (s *Server) SwapIngredientCategories(ctx context.Context, rq *connect.Request[genpb.SwapIngredientCategoriesRequest]) (*connect.Response[genpb.SwapIngredientCategoriesResponse], error) {
	aIdx, err := s.sql.GetIngredientCategorySortIndex(ctx, rq.Msg.IngredientCategoryA)
	if err != nil {
		return nil, err
	}

	bIdx, err := s.sql.GetIngredientCategorySortIndex(ctx, rq.Msg.IngredientCategoryB)
	if err != nil {
		return nil, err
	}

	err = s.sql.UpdateIngredientCategorySortIndex(ctx, gensql.UpdateIngredientCategorySortIndexParams{
		ID:        rq.Msg.IngredientCategoryA,
		SortIndex: bIdx,
	})

	if err != nil {
		return nil, err
	}

	err = s.sql.UpdateIngredientCategorySortIndex(ctx, gensql.UpdateIngredientCategorySortIndexParams{
		ID:        rq.Msg.IngredientCategoryB,
		SortIndex: aIdx,
	})

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.SwapIngredientCategoriesResponse{}), nil
}
