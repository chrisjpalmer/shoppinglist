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

func (s *Server) renderGotPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	pctx := page.NewContext(r)

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
	ingredientCounts, ingg, err := s.ingredients(ctx)
	if err != nil {
		return nil, err
	}

	var gg []page.GotItem
	for _, ing := range ingg {
		if ingredientCounts[ing.ID] == 0 && ing.WantOverrideCount == 0 {
			continue
		}
		gg = append(gg, page.GotItem{
			ID:         ing.ID,
			Ingredient: ing.Name,
			GotCount:   int(ing.GotCount),
		})
	}

	return gg, nil
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
