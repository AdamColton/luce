package httptype

import (
	"net/http"

	"github.com/adamcolton/luce/util/reflector"
)

// Common types from the http package
var (
	ResponseWriter = reflector.Type[http.ResponseWriter]()
	Request        = reflector.Type[*http.Request]()
	HandlerFunc    = reflector.Type[http.HandlerFunc]()
	Handler        = reflector.Type[http.Handler]()
	Header         = reflector.Type[http.Header]()
	Client         = reflector.Type[*http.Client]()
	Server         = reflector.Type[*http.Server]()
)
