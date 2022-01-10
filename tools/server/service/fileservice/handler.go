package fileservice

import (
	"encoding/json"
	"io"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lfile"
)

// StatusErr provides an optional interface for errors allowing errors to
// include HTTP status.
type StatusErr interface {
	Status() int
}

// Handler uses the path, file and request to generate a response.
type Handler func(path string, file lfile.File, req *service.Request) *service.Response

// ByteHandler uses the data from the file and the request to generate a
// response. Use Wrap to convert a ByteHandler to a Handler.
type ByteHandler func([]byte, *service.Request) *service.Response

// Wrap converts a ByteHandler to a Handler.
func Wrap(fn ByteHandler) Handler {
	return func(path string, f lfile.File, req *service.Request) *service.Response {
		b, err := io.ReadAll(f)
		if errResp := req.ErrCheck(err); errResp != nil {
			return errResp
		}
		return fn(b, req)
	}
}

// DirHandler uses the contents of the directory and the request to generate a
// response. Use WrapDir to convert a DirHandler to a Handler.
type DirHandler func(*lfile.DirContents, *service.Request) *service.Response

// Wrap converts a DirHandler to a Handler.
func WrapDir(fn DirHandler) Handler {
	return func(path string, f lfile.File, req *service.Request) *service.Response {
		dc, err := lfile.GetDirContents(f)
		if errResp := req.ErrCheck(err); errResp != nil {
			return errResp
		}
		dc.Path = path
		return fn(dc, req)
	}
}

var (
	// HandleBinary directly returns the binary contents of a file.
	HandleBinary = Wrap(func(b []byte, req *service.Request) *service.Response {
		return req.Response(b)
	})

	// HandleDir JSON encodes the lfile.DirContents and returns that.
	HandleDir = WrapDir(JsonDir)
)

// JsonDir fulfills DirHandler and profides the default HandleDir (via WrapDir).
// It returns a Response with the body set to a json representation of the
// lfile.DirContents.
func JsonDir(dc *lfile.DirContents, req *service.Request) *service.Response {
	resp := req.Response(nil)
	err := json.NewEncoder(resp).Encode(dc)
	if errResp := req.ErrCheck(err); errResp != nil {
		return errResp
	}
	return resp
}
