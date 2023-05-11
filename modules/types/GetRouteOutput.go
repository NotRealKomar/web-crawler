package types

import "net/http"

type GetRouteOutput func(http.ResponseWriter, *http.Request)
