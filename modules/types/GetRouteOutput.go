package types

import "net/http"

type RouteHandler func(http.ResponseWriter, *http.Request)
