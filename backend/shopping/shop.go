package shopping

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/chrisjpalmer/shoppinglist/backend/gensql"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/page"
	"github.com/chrisjpalmer/shoppinglist/backend/shopping/render"
)

func (s *Server) handleShopPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		shoppedIds, err := parseShoppedColumns(r)
		if err != nil {
			w.WriteHeader(400)
			fmt.Println("error while parsing shopped columns", err.Error())
			return
		}

		err = s.saveShoppedColumns(r.Context(), shoppedIds)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println("error updating shopped states", err.Error())
			return
		}
	}

	s.renderShopPage(w, r)
}

func (s *Server) handleShopResetPage(w http.ResponseWriter, r *http.Request) {
	err := s.sql.ResetIngredientShopped(r.Context())
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("error resetting shopped states", err.Error())
		return
	}

	s.renderShopPage(w, r)
}

func (s *Server) renderShopPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	pctx := s.pageContext(r)

	ss, err := s.shopItems(r.Context())
	if err != nil {
		w.WriteHeader(500)
		fmt.Println("error while getting shop items: ", err.Error())
		return
	}

	var opts []func(*templ.ComponentHandler)

	if q.Has(fragmentParam) {
		opts = append(opts, templ.WithFragments(q.Get(fragmentParam)))
	}

	templ.Handler(render.ShopPage(pctx, ss), opts...).ServeHTTP(w, r)
}

func (s *Server) shopItems(ctx context.Context) ([]page.ShopItem, error) {
	cats, err := s.ingredients(ctx)
	if err != nil {
		return nil, err
	}

	needCt := func(ing ingredient) int {
		return int(ing.RequiredCount) - int(ing.GotCount)
	}

	cats = filterIngredients(cats, func(ing ingredient) bool {
		return needCt(ing) > 0
	})

	var ss []page.ShopItem
	for _, cat := range cats {
		ss = append(ss, page.ShopItem{Category: cat.name})

		for _, ing := range cat.ingredients {
			ss = append(ss, page.ShopItem{
				ID:         ing.ID,
				Ingredient: ing.Name,
				NeedCount:  needCt(ing),
				Shopped:    ing.Shopped,
			})
		}
	}

	return ss, nil
}

func (s *Server) saveShoppedColumns(ctx context.Context, shoppedIds map[int64]bool) error {
	for id := range shoppedIds {
		err := s.sql.UpdateIngredientShopped(ctx, gensql.UpdateIngredientShoppedParams{
			ID:      id,
			Shopped: shoppedIds[id],
		})
		if err != nil {
			return fmt.Errorf("error updating shopped state for ingredient %d: %w", id, err)
		}
	}

	return nil
}

func parseShoppedColumns(r *http.Request) (map[int64]bool, error) {
	const maxMemory = 100000

	r.ParseMultipartForm(maxMemory)

	ids, err := shoppedRowIds(r.Form)
	if err != nil {
		return nil, fmt.Errorf("error getting shopped row ids: %w", err)
	}

	checked, err := shoppedChecked(r.Form)
	if err != nil {
		return nil, fmt.Errorf("error getting shopped checked: %w", err)
	}

	// all checked contains keys for rows which aren't checked as
	// well as those that are checked

	allChecked := make(map[int64]bool, len(ids))

	for _, id := range ids {
		allChecked[id] = checked[id]
	}

	return allChecked, nil
}

func shoppedRowIds(form url.Values) ([]int64, error) {
	const prefix = "col-shopped-row."

	idsMap := make(map[int64]bool)

	for k := range form {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		idstr := strings.TrimPrefix(k, prefix)

		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse id in shopped column hidden %q: %w", k, err)
		}

		idsMap[id] = true
	}

	ids := make([]int64, 0, len(idsMap))

	for id := range idsMap {
		ids = append(ids, id)
	}

	return ids, nil
}

func shoppedChecked(form url.Values) (map[int64]bool, error) {
	const prefix = "col-shopped."

	checked := make(map[int64]bool)

	for k := range form {
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		idstr := strings.TrimPrefix(k, prefix)

		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse id in shopped column %q: %w", k, err)
		}

		checked[id] = true
	}

	return checked, nil
}
