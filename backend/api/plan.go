package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/genpb"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
)

func (s *Server) GetPlan(ctx context.Context, rq *connect.Request[genpb.GetPlanRequest]) (*connect.Response[genpb.GetPlanResponse], error) {
	p, err := s.sql.GetPlan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		gp, err := s.createEmptyPlan(ctx)
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&genpb.GetPlanResponse{Plan: gp}), nil
	}

	if err != nil {
		return nil, err
	}

	var gp genpb.Plan
	err = unmarshalJSON(p.PlanData, &gp)
	if err != nil {
		return nil, err
	}

	sum, err := s.planSummary(ctx, &gp)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.GetPlanResponse{Plan: &gp, PlanSummary: sum}), nil
}

func (s *Server) UpdatePlan(ctx context.Context, rq *connect.Request[genpb.UpdatePlanRequest]) (*connect.Response[genpb.UpdatePlanResponse], error) {
	p, err := s.sql.GetPlan(ctx)
	if err != nil {
		return nil, err
	}

	pstr, err := marshalJSON(rq.Msg.Plan)
	if err != nil {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		if err := s.sql.CreatePlan(ctx, pstr); err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	err = s.sql.UpdatePlan(ctx, gensql.UpdatePlanParams{
		ID:       p.ID,
		PlanData: pstr,
	})

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&genpb.UpdatePlanResponse{}), nil
}

func (s *Server) createEmptyPlan(ctx context.Context) (*genpb.Plan, error) {
	p := emptyPlan()

	ps, err := marshalJSON(&p)
	if err != nil {
		return nil, err
	}

	err = s.sql.CreatePlan(ctx, ps)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *Server) planSummary(ctx context.Context, p *genpb.Plan) (*genpb.PlanSummary, error) {
	meals, err := s.sql.GetMeals(ctx)
	if err != nil {
		return nil, err
	}

	mealsmap, err := mealsMap(meals)
	if err != nil {
		return nil, err
	}

	ingg, err := s.sql.GetIngredients(ctx)
	if err != nil {
		return nil, err
	}

	smIds := selectedMealIds(p)

	ingredientCounts := make(map[int64]int32)

	for _, smId := range smIds {

		meal, ok := mealsmap[smId]
		if !ok {
			continue
		}

		for _, igref := range meal.IngredientRefs {
			ingredientCounts[igref.IngredientId] += igref.Number
		}
	}

	var igrefs []*genpb.IngredientRef
	for _, ing := range ingg {
		ct := ingredientCounts[ing.ID]
		if ct > 0 {
			igrefs = append(igrefs, &genpb.IngredientRef{IngredientId: ing.ID, Number: ct})
		}
	}

	return &genpb.PlanSummary{
		IngredientRef: igrefs,
	}, nil
}

func selectedMealIds(p *genpb.Plan) []int64 {
	var meals []int64
	for _, d := range p.Days {
		for _, cm := range d.CategoryMeals {
			if cm.MealId == 0 {
				// this is a valid case if the plan is uninitialised
				continue
			}
			meals = append(meals, cm.MealId)
		}
	}
	return meals
}

func mealsMap(meals []gensql.GetMealsRow) (map[int64]*genpb.Meal, error) {
	mealsmap := make(map[int64]*genpb.Meal)
	for _, m := range meals {
		var igrefs []*genpb.IngredientRef

		err := unmarshalJSON(m.Ingredients, &igrefs)
		if err != nil {
			return nil, err
		}

		mealsmap[m.ID] = &genpb.Meal{
			Name:           m.Name,
			Id:             m.ID,
			IngredientRefs: igrefs,
			RecipeUrl:      m.RecipeUrl,
		}
	}

	return mealsmap, nil
}

func emptyPlan() genpb.Plan {
	var days []*genpb.Day
	for range 7 {
		days = append(days, &genpb.Day{
			CategoryMeals: []*genpb.CategoryMeal{
				// 0 = lunch, 1 = dinner, 2 = snack
				{Category: genpb.Category_CATEGORY_LUNCH}, {Category: genpb.Category_CATEGORY_DINNER}, {Category: genpb.Category_CATEGORY_SNACK},
			},
		})
	}
	return genpb.Plan{
		Days: days,
	}
}

func marshalJSON(obj any) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unmarshalJSON(s string, obj any) error {
	return json.Unmarshal([]byte(s), obj)
}
