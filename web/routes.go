package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

// definition of a route for a web link
type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}


// newRouter starts web service.
func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// declare routes for web links.
var routes = []route{
	route{
		"Index",
		"GET",
		"/",
		index,
	},
	route{
		"TodoIndex",
		"GET",
		"/matches",
		matchesIndex,
	},
	route{
		"TodoIndex",
		"GET",
		"/matches/{taskId}",
		matchTaskIdIndex,
	},
	route{
		"TodoIndex",
		"GET",
		"/tsdb",
		tsdbIndex,
	},
	route{
		"TodoIndex",
		"GET",
		"/mesos",
		mesosIndex,
	},
	route{
		"TodoIndex",
		"GET",
		"/aurora",
		auroraIndex,
	},
	route{
		"TodoIndex",
		"GET",
		"/host",
		hostIndex,
	},
}
