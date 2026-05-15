package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/genpb"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
)

const (
	imageModeNone     = "none"
	imageModeInternal = "internal"
	imageModeExternal = "external"

	imageTypePreview     = "preview"
	imageTypeIngredients = "ingredients"
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

		preview, err := mapImageMeta(m.ID, imageTypePreview, m.PreviewImageMode, m.PreviewImageUrl)
		if err != nil {
			return nil, fmt.Errorf("error mapping image meta for PreviewImage: %w", err)
		}

		ingredients, err := mapImageMeta(m.ID, imageTypeIngredients, m.IngredientsImageMode, m.IngredientsImageUrl)
		if err != nil {
			return nil, fmt.Errorf("error mapping image meta for IngredientsImage: %w", err)
		}

		gmm = append(gmm, &genpb.Meal{
			Id:               m.ID,
			Name:             m.Name,
			IngredientRefs:   ig,
			RecipeUrl:        m.RecipeUrl,
			PreviewImage:     preview,
			IngredientsImage: ingredients,
		})
	}

	return connect.NewResponse(&genpb.GetMealsResponse{Meals: gmm}), nil
}

func mapImageMeta(id int64, imageType string, mode, imageUrl interface{}) (*genpb.ImageMeta, error) {
	imageUrlStr, ok := imageUrl.(string)
	if !ok {
		return nil, fmt.Errorf("error casting image url to a string")
	}

	modeStr, ok := mode.(string)
	if !ok {
		return nil, fmt.Errorf("error casting image url to a string")
	}

	return &genpb.ImageMeta{
		Mode:        mapImageMode(modeStr),
		ExternalUrl: imageUrlStr,
		InternalUrl: fmt.Sprintf("/meal/%d/image/%s", id, imageType), // leave out /api as this is part of the frontend url
	}, nil
}

func (s *Server) CreateMeal(ctx context.Context, rq *connect.Request[genpb.CreateMealRequest]) (*connect.Response[genpb.CreateMealResponse], error) {
	igstr, err := marshalJSON(rq.Msg.Meal.IngredientRefs)
	if err != nil {
		return nil, err
	}

	id, err := s.sql.CreateMeal(ctx, gensql.CreateMealParams{
		Name:                 rq.Msg.Meal.Name,
		Ingredients:          igstr,
		RecipeUrl:            rq.Msg.Meal.RecipeUrl,
		PreviewImageMode:     mapPBImageMode(rq.Msg.Meal.PreviewImage.Mode),
		PreviewImageUrl:      rq.Msg.Meal.PreviewImage.ExternalUrl,
		IngredientsImageMode: mapPBImageMode(rq.Msg.Meal.IngredientsImage.Mode),
		IngredientsImageUrl:  rq.Msg.Meal.IngredientsImage.ExternalUrl,
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
		ID:                   rq.Msg.Meal.Id,
		Name:                 rq.Msg.Meal.Name,
		Ingredients:          igstr,
		RecipeUrl:            rq.Msg.Meal.RecipeUrl,
		PreviewImageMode:     mapPBImageMode(rq.Msg.Meal.PreviewImage.Mode),
		PreviewImageUrl:      rq.Msg.Meal.PreviewImage.ExternalUrl,
		IngredientsImageMode: mapPBImageMode(rq.Msg.Meal.IngredientsImage.Mode),
		IngredientsImageUrl:  rq.Msg.Meal.IngredientsImage.ExternalUrl,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.UpdateMealResponse{}), nil
}

func mapPBImageMode(mode genpb.ImageMode) string {
	switch mode {
	case genpb.ImageMode_IM_NONE:
		return imageModeNone
	case genpb.ImageMode_IM_INTERNAL:
		return imageModeInternal
	}

	return imageModeExternal
}

func mapImageMode(mode string) genpb.ImageMode {
	switch mode {
	case imageModeNone:
		return genpb.ImageMode_IM_NONE
	case imageModeInternal:
		return genpb.ImageMode_IM_INTERNAL
	}

	return genpb.ImageMode_IM_EXTERNAL
}

func (s *Server) DeleteMeal(ctx context.Context, rq *connect.Request[genpb.DeleteMealRequest]) (*connect.Response[genpb.DeleteMealResponse], error) {
	err := s.sql.DeleteMeal(ctx, rq.Msg.MealId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.DeleteMealResponse{}), nil
}

func (s *Server) UpdateMealPreviewImageRequest(ctx context.Context, rq *connect.Request[genpb.UpdateMealImageRequest]) (*connect.Response[genpb.UpdateMealImageResponse], error) {
	err := s.sql.UpdateMealPreviewImageBytes(ctx, gensql.UpdateMealPreviewImageBytesParams{
		ID:                rq.Msg.Id,
		PreviewImageBytes: rq.Msg.ImageBytes,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.UpdateMealImageResponse{}), nil
}

func (s *Server) UpdateMealIngredientsImageRequest(ctx context.Context, rq *connect.Request[genpb.UpdateMealImageRequest]) (*connect.Response[genpb.UpdateMealImageResponse], error) {
	err := s.sql.UpdateMealIngredientsImageBytes(ctx, gensql.UpdateMealIngredientsImageBytesParams{
		ID:                    rq.Msg.Id,
		IngredientsImageBytes: rq.Msg.ImageBytes,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.UpdateMealImageResponse{}), nil
}

// handleMealPreviewImage - returns the meal preview image (non-protobuf)
func (s *Server) handleMealPreviewImage(rw http.ResponseWriter, r *http.Request) {
	s.handleMealImage(rw, r, imageTypePreview)
}

// handleMealIngredientsImage - returns the meal ingredients image (non-protobuf)
func (s *Server) handleMealIngredientsImage(rw http.ResponseWriter, r *http.Request) {
	s.handleMealImage(rw, r, imageTypeIngredients)
}

// handleMealImage - returns the meal preview/ingredients image
func (s *Server) handleMealImage(rw http.ResponseWriter, r *http.Request, imageType string) {
	ctx := r.Context()

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		rw.Write([]byte("unable to parse path parameter id as an int"))
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.mealImage(ctx, id, imageType, rw)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			rw.Write([]byte("the image could not be found in the database"))
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		rw.Write([]byte("internal error"))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

var ErrNotFound = fmt.Errorf("the data could not be found in the database")

func (s *Server) mealImage(ctx context.Context, id int64, imageType string, rw http.ResponseWriter) error {
	var (
		bd  interface{}
		err error
	)

	if imageType == imageTypePreview {
		bd, err = s.sql.GetMealPreviewImageBytes(ctx, id)
	} else {
		bd, err = s.sql.GetMealIngredientsImageBytes(ctx, id)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}

		return fmt.Errorf("error retrieving preview image from database: %w")
	}

	if err := writeImage(rw, bd); err != nil {
		return fmt.Errorf("error writing image to http request: %w", err)
	}

	return nil
}

func writeImage(rw http.ResponseWriter, bd interface{}) error {
	b, ok := bd.([]byte)
	if !ok {
		return fmt.Errorf("failed to cast binary data to bytes")
	}

	rw.Header().Set("Content-Type", "image/png") // todo make it so we can do other formats
	rw.Header().Set("Content-Length", strconv.Itoa(len(b)))

	if _, err := rw.Write(b); err != nil {
		return fmt.Errorf("failed to write data to response: %w", err)
	}

	return nil
}

func parseImageType(imageType string) (string, error) {
	switch imageType {
	case imageTypePreview, imageTypeIngredients:
		return imageType, nil
	}
	return "", fmt.Errorf("invalid image type %q: must be %q or %q", imageType, imageTypePreview, imageTypeIngredients)
}
