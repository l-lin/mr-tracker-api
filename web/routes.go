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
		"Mangas",
		"GET",
		"/mangas",
		Mangas,
	},
	Route{
		"SaveManga",
		"POST",
		"/mangas",
		SaveManga,
	},
	Route{
		"SaveMangas",
		"POST",
		"/mangas/import",
		SaveMangas,
	},
	Route{
		"Manga",
		"GET",
		"/mangas/{mangaId}",
		Manga,
	},
	Route{
		"UpdateManga",
		"PUT",
		"/mangas/{mangaId}",
		UpdateManga,
	},
	Route{
		"DeleteManga",
		"DELETE",
		"/mangas",
		DeleteMangas,
	},
	Route{
		"DeleteManga",
		"DELETE",
		"/mangas/{mangaId}",
		DeleteManga,
	},
	Route{
		"Notifications",
		"GET",
		"/notifications",
		Notifications,
	},
	Route{
		"Notification",
		"GET",
		"/notifications/{notificationId}",
		Notification,
	},
	Route{
		"DeleteNotification",
		"DELETE",
		"/notifications/{notificationId}",
		DeleteNotification,
	},
	Route{
		"NewMangas",
		"GET",
		"/newMangas",
		NewMangas,
	},
	Route{
		"AuthTest",
		"GET",
		"/authTest",
		AuthTest,
	},
}
