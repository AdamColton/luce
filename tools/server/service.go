package server

import (
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

	sr := sc.getServiceRoute(idx)
	if err := sr.GetError(); err != nil {
		panic(err.Error())
	}
	sc.routes.Add(cfg.ID)
	sc.s.serviceRoutes.Set(cfg.ID, sr)
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
			// TODO: handle error
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

func (sc *serviceConn) getServiceRoute(idx int) *serviceRoute {
	srv := sc.service
	rt := srv.Routes[idx]
	id := sc.service.Name + " " + rt.ID
	sr, found := sc.s.serviceRoutes.Get(id)
	if !found {
		var r *mux.Route
		var router = sc.s.coreserver.Router
		p := path.Join(srv.Base, rt.Path)
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
	return sr
}
