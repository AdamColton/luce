package server

import (
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lusers"
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
	routes  map[string]bool
}

func (s *Server) handleServiceSocket(netConn net.Conn) {
	defer netConn.Close()
	conn, err := service.NewConn(netConn)
	if s.Handle(err) {
		return
	}
	sc := &serviceConn{
		s:       s,
		Sender:  conn.Sender,
		respMap: make(map[uint32]chan<- *service.Response),
		routes:  make(map[string]bool),
	}
	err = bus.DefaultRegistrar.Register(conn.Listener, sc)
	s.Handle(err)

	conn.Listener.Run()

	for id := range sc.routes {
		s.serviceRoutes[id].setActive(false)
	}
}

type serviceRoute struct {
	*mux.Route
	active bool
}

func (sr *serviceRoute) setActive(active bool) {
	if sr.active == active {
		return
	}
	if active {
		ptr := reflect.
			ValueOf(sr.Route).
			Elem().
			FieldByName("buildOnly").
			Addr().
			Pointer()
		*(*bool)(unsafe.Pointer(ptr)) = false
	} else {
		sr.BuildOnly()
	}
	sr.active = active
}

func (sc *serviceConn) ResponseHandler(resp *service.Response) {
	sc.mapLock.Lock()
	ch := sc.respMap[resp.ID]
	sc.mapLock.Unlock()
	if ch == nil {
		return
	}
	ch <- resp
}

func (sc *serviceConn) RoutesHandler(routes service.Routes) {
	for _, r := range routes {
		sc.registerServiceRoute(r)
	}
}

func (sc *serviceConn) registerServiceRoute(route service.RouteConfig) {
	cvrt := sc.routeConfigToRequestConverter(route)
	h := func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Route Request: ", route.PathPrefix, route.Path)
		req := cvrt(r)
		if req == nil {
			return
		}
		ch := make(chan *service.Response)
		sc.mapLock.Lock()
		sc.respMap[req.ID] = ch
		sc.mapLock.Unlock()

		err := sc.Sender.Send(req)
		lerr.Panic(err)
		select {
		case resp := <-ch:
			if resp.Status == service.HttpRedirect {
				url := string(resp.Body)
				http.Redirect(w, r, url, resp.Status)
				break
			}
			if resp.Status > 0 {
				w.WriteHeader(resp.Status)
			}
			w.Write(resp.Body)
		case <-time.After(TimeoutDuration):
			w.WriteHeader(http.StatusRequestTimeout)
		}

		sc.mapLock.Lock()
		delete(sc.respMap, req.ID)
		sc.mapLock.Unlock()
	}

	sr := sc.s.serviceRoutes[route.ID]
	if sr == nil {
		var r *mux.Route
		if route.PathPrefix {
			r = sc.s.Router.PathPrefix(route.Path)
		} else {
			r = sc.s.Router.Path(route.Path)
		}
		if route.Method != "" {
			r = r.Methods(route.Methods()...)
		}
		if route.Host != "" {
			r = r.Host(route.Host)
		}
		sr = &serviceRoute{
			Route:  r,
			active: true,
		}
		sc.s.serviceRoutes[route.ID] = sr
	} else {
		sr.setActive(true)
	}
	sc.routes[route.ID] = true
	sr.HandlerFunc(h)
}

func (sc *serviceConn) routeConfigToRequestConverter(cfg service.RouteConfig) func(r *http.Request) *service.Request {
	var groups []string
	if cfg.Require.Group != "" {
		groups = strings.Split(cfg.Require.Group, ",")
	}

	return func(r *http.Request) *service.Request {
		var u *lusers.User
		if len(groups) > 0 || cfg.User {
			u, _ = sc.s.Users.User(r)
		}

		if !u.OneRequired(groups) {
			return nil
		}

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
			out.User = u
		}

		return out
	}
}
