package server

import (
	"net"
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/tools/server/service"
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
	s.servicesMux.Lock()
	err = bus.DefaultRegistrar.Register(conn.Listener, sc)
	s.Handle(err)
	s.servicesMux.Unlock()

	conn.Listener.Run()

	sc.routes.Each(func(id string, done *bool) {
		s.serviceRoutes.GetVal(id).setActive(false)
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
