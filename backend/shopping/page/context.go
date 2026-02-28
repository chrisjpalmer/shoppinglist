package page

import "net/http"

// Context - represents the context of the current page
type Context struct {
	// Path - the URL path of the currently viewed page
	Path string
}

// NewContext - returns a page context
func NewContext(r *http.Request) Context {
	return Context{
		Path: r.URL.Path,
	}
}
