package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"sort"

	"connectrpc.com/connect"
	"github.com/chrisjpalmer/shoppinglist/backend/gen"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
)

func (s *Server) GetPlan(ctx context.Context, rq *connect.Request[gen.GetPlanRequest]) (*connect.Response[gen.GetPlanResponse], error) {
	p, err := s.sql.GetPlan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		gp, err := s.createEmptyPlan(ctx)
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&gen.GetPlanResponse{Plan: gp}), nil
	}

	if err != nil {
		return nil, err
	}

	var gp gen.Plan
	err = unmarshalJSON(p.PlanData, &gp)
	if err != nil {
		return nil, err
	}

	sum, err := s.planSummary(ctx, &gp)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.GetPlanResponse{Plan: &gp, PlanSummary: sum}), nil
}

func (s *Server) UpdatePlan(ctx context.Context, rq *connect.Request[gen.UpdatePlanRequest]) (*connect.Response[gen.UpdatePlanResponse], error) {
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

	err = s.sql.UpdatePlan(ctx, generated.UpdatePlanParams{
		ID:       p.ID,
		PlanData: pstr,
	})

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&gen.UpdatePlanResponse{}), nil
}

func (s *Server) createEmptyPlan(ctx context.Context) (*gen.Plan, error) {
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

func (s *Server) planSummary(ctx context.Context, p *gen.Plan) (*gen.PlanSummary, error) {
	meals, err := s.sql.GetMeals(ctx)
	if err != nil {
		return nil, err
	}

	mealsmap, err := mealsMap(meals)
	if err != nil {
		return nil, err
	}

	smIds := selectedMealIds(p)

	ingredientCounts := make(map[int64]int32)

	for _, smId := range smIds {

		meal, ok := mealsmap[smId]
		if !ok {
			log.Printf("ignored selected meal id %d as it couldnt be found in the meals map", smId)
			continue
		}

		for _, igref := range meal.IngredientRefs {
			ingredientCounts[igref.IngredientId] += igref.Number
		}
	}

	var igrefs IngredientRefs
	for id, ct := range ingredientCounts {
		igrefs = append(igrefs, &gen.IngredientRef{IngredientId: id, Number: ct})
	}

	sort.Sort(igrefs)

	return &gen.PlanSummary{
		IngredientRef: igrefs,
	}, nil
}

type IngredientRefs []*gen.IngredientRef

func (x IngredientRefs) Len() int           { return len(x) }
func (x IngredientRefs) Less(i, j int) bool { return x[i].IngredientId < x[j].IngredientId }
func (x IngredientRefs) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func selectedMealIds(p *gen.Plan) []int64 {
	var meals []int64
	for _, d := range p.Days {
		for _, cm := range d.CategoryMeals {
			meals = append(meals, cm.MealId)
		}
	}
	return meals
}

func mealsMap(meals []generated.GetMealsRow) (map[int64]*gen.Meal, error) {
	mealsmap := make(map[int64]*gen.Meal)
	for _, m := range meals {
		var igrefs []*gen.IngredientRef

		err := unmarshalJSON(m.Ingredients, &igrefs)
		if err != nil {
			return nil, err
		}

		mealsmap[m.ID] = &gen.Meal{
			Name:           m.Name,
			Id:             m.ID,
			IngredientRefs: igrefs,
			RecipeUrl:      m.RecipeUrl,
		}
	}

	return mealsmap, nil
}

func emptyPlan() gen.Plan {
	var days []*gen.Day
	for range 7 {
		days = append(days, &gen.Day{
			CategoryMeals: []*gen.CategoryMeal{
				// 0 = lunch, 1 = dinner, 2 = snack
				{Category: gen.Category_CATEGORY_LUNCH}, {Category: gen.Category_CATEGORY_DINNER}, {Category: gen.Category_CATEGORY_SNACK},
			},
		})
	}
	return gen.Plan{
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
