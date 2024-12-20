package server

import (
	"io"
	"math/rand"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
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
	service *service.Service
	s       *Server
	*serialbus.Sender
	respMap lmap.Wrapper[uint32, chan<- *service.Response]
	routes  *lset.Set[string]
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
		respMap: lmap.NewSafe[uint32, chan<- *service.Response](nil),
		routes:  lset.New[string](),
	}
	err = bus.DefaultRegistrar.Register(conn.Listener, sc)
	s.Handle(err)

	conn.Listener.Run()

	sc.routes.Each(func(id string, done *bool) {
		s.serviceRoutes[id].setActive(false)
	})

	if sc.service != nil {
		s.services.Delete(sc.service.Name)
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
	ch := sc.respMap.GetVal(resp.ID)
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

func (s *Server) getServiceRoute(rc service.RouteConfig) *serviceRoute {
	sr, found := s.serviceRoutes[rc.ID]
	if !found {
		var r *mux.Route
		var router = s.coreserver.Router
		if rc.PathPrefix {
			r = router.PathPrefix(rc.Path)
		} else {
			r = router.Path(rc.Path)
		}
		if rc.Method != "" {
			r = r.Methods(rc.Methods()...)
		}
		if rc.Host != "" {
			r = r.Host(rc.Host)
		}
		sr = &serviceRoute{
			Route:  r,
			active: true,
		}
		s.serviceRoutes[rc.ID] = sr
	} else {
		sr.setActive(true)
	}
	return sr
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
		sc.respMap.Set(req.ID, ch)
		err := sc.Sender.Send(req)
		lerr.Panic(err)
		select {
		case resp := <-ch:
			if resp.Status == service.HttpRedirect {
				url := string(resp.Body)
				http.Redirect(w, r, url, resp.Status)
				break
			}

			h := w.Header()
			for key, val := range resp.Header {
				h[key] = val
			}
			if h[service.ContentType] == nil {
				ct := mime.TypeByExtension(filepath.Ext(r.URL.Path))
				if ct != "" {
					h.Set(service.ContentType, ct)
				}
			}
			if resp.Status > 0 {
				w.WriteHeader(resp.Status)
			}
			w.Write(resp.Body)
		case <-time.After(TimeoutDuration):
			w.WriteHeader(http.StatusRequestTimeout)
		}

		sc.respMap.Delete(req.ID)
	}

	sr := sc.s.getServiceRoute(route)
	sc.routes.Add(route.ID)
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
			out.Body, _ = io.ReadAll(r.Body)
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
