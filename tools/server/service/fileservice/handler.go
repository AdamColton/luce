package fileservice

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lfile"
)

type Handler func(path string, file *os.File, req *service.Request) *service.Response

type ByteHandler func([]byte, *service.Request) *service.Response

func Wrap(fn ByteHandler) Handler {
	return func(path string, f *os.File, req *service.Request) *service.Response {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return req.ResponseString(err.Error())
		}
		return fn(b, req)
	}
}

type DirHandler func(*lfile.DirContents, *service.Request) *service.Response

func WrapDir(fn DirHandler) Handler {
	return func(path string, f *os.File, req *service.Request) *service.Response {
		dc, err := lfile.GetDirContents(f)
		if err != nil {
			return req.ResponseErr(err, 500)
		}
		dc.Path = path
		return fn(dc, req)
	}
}

var (
	HandleBinary = Wrap(func(b []byte, req *service.Request) *service.Response {
		return req.Response(b)
	})
	HandleDir = WrapDir(JsonDir)
)

func JsonDir(dc *lfile.DirContents, req *service.Request) *service.Response {
	resp := req.Response(nil)
	err := json.NewEncoder(resp).Encode(dc)
	if err != nil {
		return req.ResponseErr(err, 500)
	}
	return resp
}
