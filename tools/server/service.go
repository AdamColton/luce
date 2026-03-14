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

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/gorilla/mux"
)

// So I need to remember how this works because there's some magic involved.
// But somehow this method ends up getting called.
func (sc *serviceConn) ServiceHandler(srv *service.Service) {
	sc.service = srv
	sc.s.services.Set(srv.Name, sc)
	for idx := range srv.Routes {
		sc.registerService(idx)
	}
}

func (sc *serviceConn) registerService(idx int) {
	cfg := sc.service.Routes[idx]
	cvrt := sc.routeToRequestConverter(cfg)
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

func (sc *serviceConn) routeToRequestConverter(cfg service.Route) func(r *http.Request) *service.Request {
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
