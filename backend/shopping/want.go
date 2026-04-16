package shopping

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/gen"
	"github.com/chrisjpalmer/shoppinglist/backend/generated"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/page"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
)

func (m *Server) handleWantPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		ovCt, err := parseOverrideColumns(r)
		if err != nil {
			w.WriteHeader(400)
			fmt.Println("error while parsing override columns", err.Error())
			return
		}

		err = m.saveOverrideColumns(r.Context(), ovCt)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error updating override counts", err.Error())
			return
		}
	}

	m.renderWantPage(w, r)
}

func (m *Server) renderWantPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	pctx := page.NewContext(r)

	ww, err := m.wantItems(r.Context())
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("error while getting plan summary: ", err.Error())
		return
	}

	var opts []func(*templ.ComponentHandler)

	if q.Has(fragmentParam) {
		opts = append(opts, templ.WithFragments(q.Get(fragmentParam)))
	}

	templ.Handler(render.WantPage(pctx, ww), opts...).ServeHTTP(w, r)
}

func (m *Server) saveOverrideColumns(ctx context.Context, ovCt map[int64]int64) error {
	for id, ct := range ovCt {
		err := m.sql.UpdateIngredientWantOverrideCount(ctx, generated.UpdateIngredientWantOverrideCountParams{
			ID:                id,
			WantOverrideCount: ct,
		})

		if err != nil {
			return fmt.Errorf("error updating want override count for ingredient %d: %w", id, err)
		}
	}

	return nil
}

func parseOverrideColumns(r *http.Request) (map[int64]int64, error) {
	const maxMemory = 100000

	const prefix = "col-override."

	r.ParseMultipartForm(maxMemory)

	ovCt := make(map[int64]int64, len(r.Form))

	for k, v := range r.Form {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		idstr := strings.TrimPrefix(k, prefix)

		id, ct, err := parseNumericFormValue(idstr, v)
		if err != nil {
			return nil, fmt.Errorf("could not parse override column %q: %w", k, err)
		}

		ovCt[id] = ct
	}

	return ovCt, nil
}

func parseNumericFormValue(idstr string, value []string) (id int64, ct int64, err error) {
	id, err = strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse id %q as int64: %w", idstr, err)
	}

	if len(value) == 0 {
		return 0, 0, fmt.Errorf("expected value slice not to be empty: %w", err)
	}

	v := value[0]

	ct, err = strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse value %q as int64: %w", v, err)
	}

	return id, ct, nil
}

func (s *Server) wantItems(ctx context.Context) ([]page.WantItem, error) {
	ingredientCounts, ingg, err := s.ingredients(ctx)
	if err != nil {
		return nil, err
	}

	var ww []page.WantItem
	for _, ing := range ingg {
		ct := ingredientCounts[ing.ID]
		ww = append(ww, page.WantItem{
			ID:            ing.ID,
			Ingredient:    ing.Name,
			Count:         int(ct),
			OverrideCount: int(ing.WantOverrideCount),
		})
	}

	return ww, nil
}

func (s *Server) ingredients(ctx context.Context) (map[int64]int32, []generated.Ingredient, error) {
	plan, err := s.plan(ctx)
	if err != nil {
		return nil, nil, err
	}

	meals, err := s.sql.GetMeals(ctx)
	if err != nil {
		return nil, nil, err
	}

	mealsmap, err := mealsMap(meals)
	if err != nil {
		return nil, nil, err
	}

	ingg, err := s.sql.GetIngredients(ctx)
	if err != nil {
		return nil, nil, err
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

	return ingredientCounts, ingg, nil
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
