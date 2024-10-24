package service

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"net/url"

	"github.com/adamcolton/luce/lhttp"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/adamcolton/luce/util/lusers"
)

var tm = type32.NewTypeMap()
var OS lfile.FSFileReader = lfile.OSRepository{}

// Register types with both gob and the typemap.
func Register(zeroValues ...type32.TypeIDer32) {
	for _, z := range zeroValues {
		gob.Register(z)
	}
	tm.RegisterType32s(zeroValues...)
}

func init() {
	Register(
		(Routes)(nil),
		(*Request)(nil),
		(*Response)(nil),
		SocketOpened{},
		SocketClose{},
		SocketMessage{},
	)
}

// Request represents a user request that the luce Server is relaying to the
// service.
type Request struct {
	ID          uint32
	RouteConfig string
	Path        string
	Method      string
	PathVars    map[string]string
	Query       map[string]string
	Form        url.Values
	Body        []byte
	User        *lusers.User
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (*Request) TypeID32() uint32 {
	return 161709784
}

// Response to the Request.
func (r *Request) Response(body []byte) *Response {
	return &Response{
		ID:     r.ID,
		Body:   body,
		Status: http.StatusOK,
	}
}

func (r *Request) ServeFile(path string) (resp *Response) {
	file, err := OS.ReadFile(path)
	resp = r.ErrCheck(err)
	if resp == nil {
		resp = r.Response(file)
	}
	return
}

// SerializeResponse uses the provided Serializer and data to create a Response
// to the Request.
func (r *Request) SerializeResponse(s serial.Serializer, data any, buf []byte) (*Response, error) {
	body, err := s.Serialize(data, buf)
	if err != nil {
		return r.ResponseErr(err, 500), err
	}
	return r.Response(body), nil

}

// ResponseTemplate populates the body of the response using the provided
// template and data. If there is a template error, that will populate the
// body.
func (r *Request) ResponseTemplate(name string, t luceio.TemplateExecutor, data any) *Response {
	buf := bytes.NewBuffer(nil)
	var err error
	out := r.Response(nil)

	if name == "" {
		err = t.Execute(buf, data)
	} else {
		err = t.ExecuteTemplate(buf, name, data)
	}

	if err != nil {
		out.Body = []byte(err.Error())
		out.Status = http.StatusInternalServerError
	} else {
		out.Body = buf.Bytes()
	}

	return out
}

// ResponseString to the Request.
func (r *Request) ResponseString(body string) *Response {
	return r.Response([]byte(body))
}

// ResponseErr sets the response body to the error and sets the status.
func (r *Request) ResponseErr(err error, status int) *Response {
	resp := r.ResponseString(err.Error())
	resp.Status = status
	return resp
}

const HttpRedirect = 302

func (r *Request) Redirect(url string) *Response {
	resp := r.ResponseString(url)
	resp.Status = HttpRedirect
	return resp
}

// ResponseErr sets the response body to the error and sets the status.
func (r *Request) ErrCheck(err error) *Response {
	s := lhttp.ErrStatus(err)
	if s == 0 {
		return nil
	}
	return r.
		ResponseString(err.Error()).
		SetStatus(s)
}

// Response to a request. The ID is the same as the ID is taken from the
// request.
type Response struct {
	ID     uint32
	Body   []byte
	Status int
	Header http.Header
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (*Response) TypeID32() uint32 {
	return 370114636
}

// Set Response HTTP Status Code
func (r *Response) SetStatus(status int) *Response {
	r.Status = status
	return r
}

// ErrCheck will set the response to the error body only if err is not nil. In
// this case it will check if the error fulfills lhttp.StatusErr and use that
// for the status, other wise it will set the status to
// StatusInternalServerError.
func (r *Response) ErrCheck(err error) (notNil bool) {
	notNil = err != nil
	if notNil {
		r.Body = []byte(err.Error())
		r.Status = http.StatusInternalServerError
		if s, ok := err.(lhttp.StatusErr); ok {
			r.Status = s.Status()
		}
	}
	return
}

const ContentType = "Content-Type"

func (r *Response) ContentType(val string) *Response {
	return r.SetHeader(ContentType, val)
}

func (r *Response) SetHeader(key, val string) *Response {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Set(key, val)
	return r
}

// Write p to the body. Fulfills io.Writer.
func (r *Response) Write(p []byte) (n int, err error) {
	if r.Body == nil {
		r.Body = p
	} else {
		r.Body = append(r.Body, p...)
	}
	return len(p), nil
}

type SocketOpened struct {
	ID uint32
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketOpened) TypeID32() uint32 {
	return 1046109042
}

type SocketClose struct {
	ID uint32
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketClose) TypeID32() uint32 {
	return 3196974518
}

type SocketMessage struct {
	ID   uint32
	Body []byte
}

// TypeID32 fulfill TypeIDer32. The ID was choosen at random.
func (SocketMessage) TypeID32() uint32 {
	return 3196974518
}
