package server

import (
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/gorilla/mux"
)

// This is wired up using magic. Because it ends in Handler, it's wired by type
func (sc *serviceConn) ServiceHandler(srv *service.Service) {
	sc.service = srv
	sc.s.services.Set(srv.Name, sc)
	for idx := range srv.Routes {
		sc.registerService(idx)
	}
}

// SocketMessageHandler is wired up the same way as ServiceHandler: a
// registered service sends a SocketMessage to push Body out over the
// client websocket it opened, identified by msg.ID.
func (sc *serviceConn) SocketMessageHandler(msg service.SocketMessage) {
	ws := sc.s.byId.GetVal(msg.ID)
	if ws == nil {
		return
	}
	ws.Write(msg.Body)
}

// SocketCloseHandler is wired up the same way as ServiceHandler: a
// registered service sends a SocketClose to close the client websocket it
// opened, identified by msg.ID, and stop tracking it.
func (sc *serviceConn) SocketCloseHandler(msg service.SocketClose) {
	sc.sockets.Remove(msg.ID)
	sc.s.closeSocket(msg.ID)
}

// closeSocket closes the client websocket id names and stops tracking it.
// Closing it ends the proxy's read loop, which tidies up the rest.
func (s *Server) closeSocket(id uint32) {
	ws := s.WebSocket.byId.GetVal(id)
	if ws == nil {
		return
	}
	if c, ok := ws.MessageReaderWriter.(io.Closer); ok {
		c.Close()
	}
	s.WebSocket.byId.Delete(id)
}

func (sc *serviceConn) registerService(idx int) {
	cfg := sc.service.Routes[idx]
	cvrt := sc.routeToRequestConverter(cfg)
	h := func(w http.ResponseWriter, r *http.Request) {
		req := cvrt(r, w)
		if req == nil {
			return
		}
		if req.WebSocket > 0 {
			sc.websocket(w, r, req)
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
				// == projects.Code.luce.server ==
				// [ ] Content-Type guessed from the request path, not the body
				//  A service's Base usually ends in ".service", and
				//  mime.TypeByExtension(".service") is "text/x-dbus-service", so a
				//  page served at the bare Base is sent with that and the browser
				//  renders the HTML as plain text. Serving it one path down, as
				//  /stt.service/home, is the way around it for now. The guess
				//  should probably only apply to a route that serves files, or a
				//  handler should be able to say "no Content-Type" and leave the
				//  browser to sniff.
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

	id, sr := sc.getServiceRoute(idx)
	if err := sr.GetError(); err != nil {
		panic(err.Error())
	}
	sc.routes.Add(id)
	sr.HandlerFunc(h)
}

// websocket proxies the client websocket req.WebSocket names: it sends req
// to the registered service as a normal request, then reads frames off the
// upgraded socket and forwards each as a SocketMessage until the socket
// closes, at which point it sends a SocketClose.
func (sc *serviceConn) websocket(w http.ResponseWriter, r *http.Request, req *service.Request) {
	// Tracked from here on, so that if the service disconnects while this
	// socket is open the teardown knows to close it.
	sc.sockets.Add(req.WebSocket)
	defer sc.sockets.Remove(req.WebSocket)

	ch := make(chan *service.Response)
	sc.respMap.Set(req.ID, ch)
	if sc.Sender.Send(req) != nil {
		sc.s.closeSocket(req.WebSocket)
		sc.respMap.Delete(req.ID)
		return
	}
	select {
	case <-ch:
	case <-time.After(TimeoutDuration):
	}
	sc.respMap.Delete(req.ID)
	// nothing is done with the reponse...

	ws := sc.s.WebSocket.byId.GetVal(req.WebSocket)
	if ws == nil {
		// Closed while the open was in flight, which is what the teardown
		// does to every socket when the service disconnects.
		return
	}
	from := make(chan []byte)
	stop := false
	go (iobus.Config{
		CloseOnEOF: true,
	}).Reader(ws, from, nil, &stop)
	for body := range from {
		msg := service.SocketMessage{
			Body: body,
			ID:   req.WebSocket,
		}
		// A send that fails means the service is gone. Its sockets are closed
		// by the teardown, so this only has to stop proxying.
		if sc.Sender.Send(msg) != nil {
			return
		}
	}

	cls := service.SocketClose{
		ID: req.WebSocket,
	}
	if sc.Sender.Send(cls) != nil {
		return
	}
	sc.s.WebSocket.byId.Delete(req.WebSocket)
}

func (sc *serviceConn) routeToRequestConverter(cfg service.Route) func(r *http.Request, w http.ResponseWriter) *service.Request {
	var groups []string
	if cfg.Require.Group != "" {
		groups = strings.Split(cfg.Require.Group, ",")
	}

	return func(r *http.Request, w http.ResponseWriter) *service.Request {
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
			// == projects.Code.luce.server ==
			// [ ] handle ParseForm's error in routeToRequestConverter
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

		if cfg.WebSocket {
			out.WebSocket = sc.s.WebSocket.counter.Add(1)
			ws, err := sc.s.WebSocket.Upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
			sc.s.WebSocket.byId.Set(out.WebSocket, lhttp.NewSocket(ws))
		}

		if cfg.User {
			out.User = u
		}

		return out
	}
}

// == projects.Code.luce.server ==
// [ ] remove stale service routes from the router
//  A route whose path, host or method changes between runs of its service
//  gets a new mux.Route, and the old one is left in the router as BuildOnly
//  until the server restarts. It never matches, so it's harmless, but each
//  change adds one. gorilla/mux has no way to remove a route, so this would
//  need the router rebuilt or a router of our own in front of it.

// getServiceRoute returns the route at idx, creating it the first time the
// service registers it and reactivating it after a reconnect. The id it
// returns is the route's key in serviceRoutes.
func (sc *serviceConn) getServiceRoute(idx int) (string, *serviceRoute) {
	srv := sc.service
	rt := srv.Routes[idx]
	p := path.Join(srv.Base, rt.Path)
	// The key holds everything that shapes the mux.Route, so a route that
	// comes back with a different path, host or method gets a new one rather
	// than reactivating the old one where it no longer belongs. The name
	// scopes it to the service.
	id := fmt.Sprintf("%s %s %s (%s) %s %t", srv.Name, rt.ID, srv.Host, rt.Method, p, rt.PathPrefix)
	sr, found := sc.s.serviceRoutes.Get(id)
	if !found {
		var r *mux.Route
		var router = sc.s.coreserver.Router
		if rt.PathPrefix {
			r = router.PathPrefix(p)
		} else {
			r = router.Path(p)
		}
		if rt.Method != "" {
			r = r.Methods(rt.Methods()...)
		}
		if srv.Host != "" {
			r = r.Host(srv.Host)
		}
		sr = &serviceRoute{
			Route:  r,
			active: true,
		}
		sc.s.serviceRoutes.Set(id, sr)
	} else {
		sr.setActive(true)
	}
	return id, sr
}
