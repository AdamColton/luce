// Package httptype provides reflect.Type values for common net/http types.
package httptype

import (
	"net/http"

	"github.com/adamcolton/luce/util/reflector"
)

// The reflect.Type of common types from net/http. Request, Client and Server
// are the pointer types, as they are normally used.
var (
	ResponseWriter = reflector.Type[http.ResponseWriter]()
	Request        = reflector.Type[*http.Request]()
	HandlerFunc    = reflector.Type[http.HandlerFunc]()
	Handler        = reflector.Type[http.Handler]()
	Header         = reflector.Type[http.Header]()
	Client         = reflector.Type[*http.Client]()
	Server         = reflector.Type[*http.Server]()
)
