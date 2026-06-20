package shopping

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/genpb"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/page"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
)

func (m *Server) handleWantPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		want, err := parseWantColumns(r)
		if err != nil {
			w.WriteHeader(400)
			fmt.Println("error while parsing override columns", err.Error())
			return
		}

		err = m.saveWantColumns(r.Context(), want)
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

	pctx := m.pageContext(r)

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

func (m *Server) saveWantColumns(ctx context.Context, want map[int64]*wantColumns) error {
	for id, w := range want {
		err := m.sql.UpdateIngredientCounts(ctx, gensql.UpdateIngredientCountsParams{
			ID:                id,
			WantOverrideCount: w.WantOverrideCount,
			MinCount:          w.MinCount,
			MaxCount:          w.MaxCount,
		})

		if err != nil {
			return fmt.Errorf("error updating counts for ingredient %d: %w", id, err)
		}
	}

	return nil
}

type wantColumns struct {
	MinCount          int64
	MaxCount          int64
	WantOverrideCount int64
}

func parseWantColumns(r *http.Request) (map[int64]*wantColumns, error) {
	const maxMemory = 100000

	const (
		ovr = "col-override"
		min = "col-min"
		max = "col-max"
	)

	r.ParseMultipartForm(maxMemory)

	want := make(map[int64]*wantColumns, len(r.Form))

	// remove all query paramters from the form data
	for q := range r.URL.Query() {
		delete(r.Form, q)
	}

	for k, v := range r.Form {
		name, id, err := parseFormKey(k)
		if err != nil {
			return nil, fmt.Errorf("error parsing id for key %s: %w", k, err)
		}

		if len(v) == 0 {
			return nil, fmt.Errorf("empty value for key: %s", k)
		}

		ct, err := strconv.ParseInt(v[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing value for key %s: %w", k, err)
		}

		if _, ok := want[id]; !ok {
			want[id] = &wantColumns{}
		}

		switch name {
		case ovr:
			want[id].WantOverrideCount = ct
		case min:
			want[id].MinCount = ct
		case max:
			want[id].MaxCount = ct
		}
	}

	return want, nil
}

var idreg = regexp.MustCompile(`^(\S+)\.([0-9]+)$`)

func parseFormKey(key string) (string, int64, error) {
	m := idreg.FindStringSubmatch(key)
	if len(m) != 3 {
		return "", 0, fmt.Errorf("unable to parse key as id: %s", key)
	}

	name := m[1]

	sid := m[2]

	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("could not parse %q as int64: %w", sid, err)
	}

	return name, id, nil
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
	cats, err := s.ingredients(ctx)
	if err != nil {
		return nil, err
	}

	var ww []page.WantItem
	for _, cat := range cats {
		ww = append(ww, page.WantItem{Category: cat.name})

		for _, ing := range cat.ingredients {
			ww = append(ww, page.WantItem{
				ID:            ing.ID,
				Ingredient:    ing.Name,
				Required:      int(ing.PlannedCount),
				MinCount:      int(ing.MinCount),
				MaxCount:      int(ing.MaxCount),
				OverrideCount: int(ing.WantOverrideCount),
				Total:         int(ing.RequiredCount),
			})
		}
	}

	return ww, nil
}

type category struct {
	name        string
	ingredients []ingredient
}

type ingredient struct {
	ID                   int64
	Name                 string
	IngredientCategoryID int64
	PlannedCount         int64
	MinCount             int64
	MaxCount             int64
	WantOverrideCount    int64
	RequiredCount        int64
	GotCount             int64
	Shopped              bool
}

const unknownCategoryID = -1

// ingredients - returns the ingredients indexed by category id
func (s *Server) ingredients(ctx context.Context) ([]category, error) {
	igs, err := s.sql.GetIngredients(ctx)
	if err != nil {
		return nil, err
	}

	cats, err := s.sql.GetIngredientCategories(ctx)
	if err != nil {
		return nil, err
	}

	reqCt, err := s.requiredCounts(ctx)
	if err != nil {
		return nil, err
	}

	igsMap := make(map[int64][]ingredient, len(igs))

	for _, ig := range igs {
		catID := ig.IngredientCategoryID

		if !categoryExists(cats, ig.IngredientCategoryID) {
			catID = unknownCategoryID
		}

		req := reqCt[ig.ID]

		if ig.MinCount != 0 {
			req = max(ig.MinCount, req)
		}

		if ig.MaxCount != 0 {
			req = min(ig.MaxCount, req)
		}

		if ig.WantOverrideCount != 0 {
			req = ig.WantOverrideCount
		}

		igsMap[catID] = append(igsMap[catID], ingredient{
			ID:                   ig.ID,
			Name:                 ig.Name,
			IngredientCategoryID: catID,
			PlannedCount:         reqCt[ig.ID],
			MinCount:             ig.MinCount,
			MaxCount:             ig.MaxCount,
			WantOverrideCount:    ig.WantOverrideCount,
			RequiredCount:        req,
			GotCount:             ig.GotCount,
			Shopped:              ig.Shopped,
		})
	}

	outCats := make([]category, 0, len(cats))

	for _, cat := range cats {
		if len(igsMap[cat.ID]) == 0 {
			continue
		}

		outCats = append(outCats, category{
			name:        cat.Name,
			ingredients: igsMap[cat.ID],
		})
	}

	if len(igsMap[unknownCategoryID]) != 0 {
		outCats = append(outCats, category{
			name:        "Miscellaneous",
			ingredients: igsMap[unknownCategoryID],
		})
	}

	return outCats, nil
}

func categoryExists(cats []gensql.IngredientCategory, catID int64) bool {
	for _, cat := range cats {
		if cat.ID == catID {
			return true
		}
	}

	return false
}

func (s *Server) requiredCounts(ctx context.Context) (map[int64]int64, error) {
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

	smIds := selectedMealIds(plan)

	ingredientCounts := make(map[int64]int64)

	for _, smId := range smIds {
		meal, ok := mealsmap[smId]
		if !ok {
			continue
		}

		for _, igref := range meal.IngredientRefs {
			ingredientCounts[igref.IngredientId] += int64(igref.Number)
		}
	}

	return ingredientCounts, nil
}

func (s *Server) plan(ctx context.Context) (*genpb.Plan, error) {
	p, err := s.sql.GetPlan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return emptyPlan(), nil
		}

		return nil, err
	}

	var plan genpb.Plan
	err = unmarshalJSON(p.PlanData, &plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func emptyPlan() *genpb.Plan {
	var days []*genpb.Day
	for range 7 {
		days = append(days, &genpb.Day{
			CategoryMeals: []*genpb.CategoryMeal{
				// 0 = lunch, 1 = dinner, 2 = snack
				{Category: genpb.Category_CATEGORY_LUNCH}, {Category: genpb.Category_CATEGORY_DINNER}, {Category: genpb.Category_CATEGORY_SNACK},
			},
		})
	}
	return &genpb.Plan{
		Days: days,
	}
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

func unmarshalJSON(s string, obj any) error {
	return json.Unmarshal([]byte(s), obj)
}
