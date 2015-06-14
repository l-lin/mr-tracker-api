package web

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Availables routes
var routes = Routes{
	Route{
		"AuthTest",
		"GET",
		"/authTest",
		AuthTest,
	},
}
