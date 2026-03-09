package shopping

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/gen"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/page"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
)

func (m *Server) handleWantPage(w http.ResponseWriter, r *http.Request) {
	pctx := page.NewContext(r)

	ww, err := m.wantItems(r.Context())
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("error while getting plan summary: ", err.Error())
		return
	}

	templ.Handler(render.WantPage(pctx, ww)).ServeHTTP(w, r)
}

func (s *Server) wantItems(ctx context.Context) ([]page.WantItem, error) {
	plan, err := s.plan(ctx)
	if err != nil {
		return nil, err
	}

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

	smIds := selectedMealIds(plan)

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

	var ww []page.WantItem
	for _, ing := range ingg {
		ct := ingredientCounts[ing.ID]
		if ct > 0 {
			ww = append(ww, page.WantItem{
				Ingredient: ing.Name,
				Count:      ct,
			})
		}
	}

	return ww, nil
}

func (s *Server) plan(ctx context.Context) (*gen.Plan, error) {
	p, err := s.sql.GetPlan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return emptyPlan(), nil
		}

		return nil, err
	}

	var plan gen.Plan
	err = unmarshalJSON(p.PlanData, &plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func emptyPlan() *gen.Plan {
	var days []*gen.Day
	for range 7 {
		days = append(days, &gen.Day{
			CategoryMeals: []*gen.CategoryMeal{
				// 0 = lunch, 1 = dinner, 2 = snack
				{Category: gen.Category_CATEGORY_LUNCH}, {Category: gen.Category_CATEGORY_DINNER}, {Category: gen.Category_CATEGORY_SNACK},
			},
		})
	}
	return &gen.Plan{
		Days: days,
	}
}

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

func unmarshalJSON(s string, obj any) error {
	return json.Unmarshal([]byte(s), obj)
}
