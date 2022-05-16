package server

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/gorilla/mux"
)

func (s *Server) RunServiceSocket() error {
	sck := unixsocket.New(s.ServiceSocket, s.handleServiceSocket)
	return sck.Run()
}

type serviceConn struct {
	s *Server
	*serialbus.Sender
	respMap map[uint32]chan<- *service.Response
	mapLock sync.Mutex
}

func (s *Server) handleServiceSocket(netConn net.Conn) {
	defer netConn.Close()
	conn, err := service.NewConn(netConn)
	if err != nil {
		fmt.Println(err)
		return
	}
	sc := &serviceConn{
		s:       s,
		Sender:  conn.Sender,
		respMap: make(map[uint32]chan<- *service.Response),
	}
	bus.RegisterHandlerType(conn.Listener, sc)

	conn.Listener.Run()
}

func (sc *serviceConn) HandleResponse(resp *service.Response) {
	sc.mapLock.Lock()
	ch := sc.respMap[resp.ID]
	sc.mapLock.Unlock()
	if ch == nil {
		return
	}
	ch <- resp
}

func (sc *serviceConn) HandleRoutes(routes service.Routes) {
	for _, r := range routes {
		sc.registerServiceRoute(r)
	}
}

func (sc *serviceConn) registerServiceRoute(route service.RouteConfig) {
	cvrt := sc.routeConfigToRequestConverter(route)
	h := func(w http.ResponseWriter, r *http.Request) {
		req := cvrt(r)
		ch := make(chan *service.Response)
		sc.mapLock.Lock()
		sc.respMap[req.ID] = ch
		sc.mapLock.Unlock()

		err := sc.Sender.Send(req)
		lerr.Panic(err)
		select {
		case resp := <-ch:
			if resp.Status > 0 {
				w.WriteHeader(resp.Status)
			}
			w.Write(resp.Body)
		case <-time.After(time.Second * 5):
			w.WriteHeader(http.StatusRequestTimeout)
		}

		sc.mapLock.Lock()
		delete(sc.respMap, req.ID)
		sc.mapLock.Unlock()
	}

	r := sc.s.serviceRoutes[route.ID]
	if r == nil {
		if route.PathPrefix {
			r = sc.s.Router.PathPrefix(route.Path)
		} else {
			r = sc.s.Router.Path(route.Path)
		}
		if route.Method != "" {
			r = r.Methods(route.Methods()...)
		}
		sc.s.serviceRoutes[route.ID] = r
	}
	r.HandlerFunc(h)
}

func (sc *serviceConn) routeConfigToRequestConverter(cfg service.RouteConfig) func(r *http.Request) *service.Request {
	return func(r *http.Request) *service.Request {
		out := &service.Request{
			Path:        r.URL.Path,
			RouteConfig: cfg.ID,
			ID:          rand.Uint32(),
			Method:      r.Method,
		}

		if cfg.Body {
			out.Body, _ = ioutil.ReadAll(r.Body)
		}

		if cfg.Form {
			r.ParseForm()
			out.Form = r.Form
		}

		if cfg.PathVars {
			out.PathVars = mux.Vars(r)
		}

		if cfg.Query {
			q := r.URL.Query()
			if ln := len(q); ln > 0 {
				out.Query = make(map[string]string, ln)
				for k, v := range q {
					out.Query[k] = v[0]
				}
			}

		}

		if cfg.User {
			u, _ := sc.s.Users.User(r)
			out.User = u
		}

		return out
	}
}
