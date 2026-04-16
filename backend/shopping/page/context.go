package page

import "net/http"

// Context - represents the context of the current page
type Context struct {
	// Path - the URL path of the currently viewed page
	Path string
	// PlanningSiteURL - the URL path to the planning site
	PlanningSiteURL string
}

// NewContext - returns a page context
func NewContext(r *http.Request, planningSiteURL string) Context {
	return Context{
		PlanningSiteURL: planningSiteURL,
		Path:            r.URL.Path,
	}
}

// WantItem - an item to display on the want page
type WantItem struct {
	ID            int64
	Ingredient    string
	Required      int
	OverrideCount int
	Total         int
}

// GotItem - an item to display on the got page
type GotItem struct {
	ID         int64
	Ingredient string
	GotCount   int
}
