package web

import "github.com/gorilla/mux"

// Returns the routers for novels, feeds and notifications
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Headers("Content-Type", "application/json", "X-Requested-With", "XMLHttpRequest")

	for _, route := range routes {
		router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(WrapWithCheckAuth(route.HandlerFunc))
	}

	return router
}

// Returns a router for signing in Google account
func NewSignInRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Headers("Content-Type", "application/json", "X-Requested-With", "XMLHttpRequest")
	route := Route{
		"SignIn",
		"GET",
		"/signin",
		SignIn,
	}
	router.Methods(route.Method).
	Path(route.Pattern).
	Name(route.Name).
	Handler(route.HandlerFunc)

	return router
}
