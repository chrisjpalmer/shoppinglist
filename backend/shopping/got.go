package shopping

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/page"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
)

func (s *Server) handleGotPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		gotCt, err := parseGotColumns(r)
		if err != nil {
			w.WriteHeader(400)
			fmt.Println("error while parsing got columns", err.Error())
			return
		}

		err = s.saveGotColumns(r.Context(), gotCt)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error updating got counts", err.Error())
			return
		}
	}

	s.renderGotPage(w, r)
}

func (s *Server) handleGotResetPage(w http.ResponseWriter, r *http.Request) {
	err := s.sql.ResetIngredientGotCount(r.Context())
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("error resetting got counts", err.Error())
		return
	}

	s.renderGotPage(w, r)
}

func (s *Server) renderGotPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	pctx := s.pageContext(r)

	gg, err := s.gotItems(r.Context())
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("error while getting plan summary: ", err.Error())
		return
	}

	var opts []func(*templ.ComponentHandler)

	if q.Has(fragmentParam) {
		opts = append(opts, templ.WithFragments(q.Get(fragmentParam)))
	}

	templ.Handler(render.GotPage(pctx, gg), opts...).ServeHTTP(w, r)
}

func (s *Server) gotItems(ctx context.Context) ([]page.GotItem, error) {
	cats, err := s.ingredients(ctx)
	if err != nil {
		return nil, err
	}

	cats = filterIngredients(cats, func(ing ingredient) bool {
		return !(ing.RequiredCount == 0 && ing.WantOverrideCount == 0)
	})

	var gg []page.GotItem
	for _, cat := range cats {
		gg = append(gg, page.GotItem{Category: cat.name})

		for _, ing := range cat.ingredients {
			gg = append(gg, page.GotItem{
				ID:         ing.ID,
				Ingredient: ing.Name,
				GotCount:   int(ing.GotCount),
			})
		}
	}

	return gg, nil
}

func filterIngredients(cats []category, filter func(ingredient) bool) []category {
	outCats := make([]category, 0, len(cats))

	for _, cat := range cats {
		outIgs := make([]ingredient, 0, len(cat.ingredients))

		for _, ing := range cat.ingredients {
			if filter(ing) {
				outIgs = append(outIgs, ing)
			}
		}

		if len(outIgs) == 0 {
			continue
		}

		outCats = append(outCats, category{
			name:        cat.name,
			ingredients: outIgs,
		})
	}

	return outCats
}

func (s *Server) saveGotColumns(ctx context.Context, gotCt map[int64]int64) error {
	for id, ct := range gotCt {
		err := s.sql.UpdateIngredientGotCount(ctx, gensql.UpdateIngredientGotCountParams{
			ID:       id,
			GotCount: ct,
		})

		if err != nil {
			return fmt.Errorf("error updating got count for ingredient %d: %w", id, err)
		}
	}

	return nil
}

func parseGotColumns(r *http.Request) (map[int64]int64, error) {
	const maxMemory = 100000

	const prefix = "col-got."

	r.ParseMultipartForm(maxMemory)

	gotCt := make(map[int64]int64, len(r.Form))

	for k, v := range r.Form {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		idstr := strings.TrimPrefix(k, prefix)

		id, ct, err := parseNumericFormValue(idstr, v)
		if err != nil {
			return nil, fmt.Errorf("could not parse got column %q: %w", k, err)
		}

		gotCt[id] = ct
	}

	return gotCt, nil
}
