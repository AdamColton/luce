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

// RunServiceSocket listens on the configured service socket, blocking until
// it stops. External services connect here to register their routes and
// handle proxied requests.
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
	// sockets is every client websocket this service is proxying. A proxy
	// runs on its request's goroutine while the teardown runs on this
	// connection's, so it has to be the safe variant.
	sockets *lset.Set[uint32]
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
		sockets: lset.Safe[uint32](),
	}
	s.servicesMux.Lock()
	err = bus.DefaultRegistrar.Register(conn.Listener, sc)
	s.Handle(err)
	s.servicesMux.Unlock()

	conn.Listener.Run()

	sc.routes.Each(func(id string, done *bool) {
		s.serviceRoutes.GetVal(id).setActive(false)
	})

	// Close the client websockets this service was proxying. Nothing else
	// would: each proxy is blocked reading its client, and the client holds a
	// connection that looks alive with no way to learn the service behind it
	// has gone. Closing them is what makes clients reconnect, which is how
	// they find the service again once it is back.
	sc.sockets.Each(func(id uint32, done *bool) {
		s.closeSocket(id)
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
