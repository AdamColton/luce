package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/lhttp"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/gorilla/websocket"
)

// WebSocket can perform two different midware operations. It can attach a
// *websocket.Conn directly or it can attach a set of []byte channels connected
// to a websocket.
type WebSocket struct {
	Upgrader       websocket.Upgrader
	ToBuf, FromBuf int
	lhttp.ErrHandler
}

// // NewWebSocket with default buffer sizes.
func NewWebSocket() WebSocket {
	return WebSocket{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// // ChannelInitilizer creates an Initilizer that will abstract the websocket to
// // a pair of []byte channels.
func (ws WebSocket) Initilizer(to, from, conn string) Initilizer {
	return webSocketInitilizer{
		WebSocket: ws,
		to:        to,
		from:      from,
		conn:      conn,
	}
}

type webSocketInitilizer struct {
	WebSocket
	to, from, conn string
}

var (
	inChanType = reflector.Type[chan<- []byte]()
	checkSend  = filter.IsType(inChanType).
			Check(filter.TypeErr("Invalid Websocket 'To' type: %s"))

	outChanType = reflector.Type[<-chan []byte]()
	checkRecv   = filter.IsType(outChanType).
			Check(filter.TypeErr("Invalid Websocket 'From' type: %s"))

	connType  = reflector.Type[*websocket.Conn]()
	checkConn = filter.IsType(connType).
			Check(filter.TypeErr("Invalid Websocket 'Conn' type: %s"))
)

// Initilize checks the "to" and "from" fields to validate the names and types.
// It will panic if they are not.
func (i webSocketInitilizer) Initilize(fieldType reflect.Type) DataInserter {
	var hasToField, hasFromField, hasConnField bool
	var toField, fromField, connField reflect.StructField

	if fieldName(i.to) {
		toField, hasToField = fieldType.FieldByName(i.to)
	}
	if fieldName(i.from) {
		fromField, hasFromField = fieldType.FieldByName(i.from)
	}
	if fieldName(i.conn) {
		connField, hasConnField = fieldType.FieldByName(i.conn)
	}

	if !hasToField && !hasFromField && !hasConnField {
		return nil
	}
	di := webSocketDataInserter{
		WebSocket: i.WebSocket,
	}
	if hasToField {
		checkSend.Panic(toField.Type)
		di.to = toField.Index
	}
	if hasFromField {
		checkRecv.Panic(fromField.Type)
		di.from = fromField.Index
	}
	if hasConnField {
		checkConn.Panic(connField.Type)
		di.conn = connField.Index
	}
	return di
}

type webSocketDataInserter struct {
	WebSocket
	to, from, conn []int
}

// Insert sets the "to" and "from" channels on the data field.
func (di webSocketDataInserter) Insert(w http.ResponseWriter, r *http.Request, data reflect.Value) (func(), error) {
	fbi := data.Elem().FieldByIndex
	set := func(idx []int, v any) {
		fbi(idx).Set(reflect.ValueOf(v))
	}

	conn, err := di.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	if di.conn != nil {
		set(di.conn, conn)
	}

	var to chan []byte
	if di.to != nil || di.from != nil {
		sw := lhttp.NewSocket(conn)
		if di.to != nil {
			to = make(chan []byte, di.ToBuf)
			set(di.to, to)
			go iobus.Writer(sw, to, nil)
		}
		if di.from != nil {
			from := make(chan []byte, di.FromBuf)
			set(di.from, from)

			go (iobus.Config{
				CloseOnEOF: true,
			}).Reader(sw, from, nil)
		}
	}

	callbackClose := func() {
		conn.Close()
		defer func() {
			// ignore double close error
			recover()
		}()
		if to != nil {
			close(to)
		}
	}
	return callbackClose, nil
}
