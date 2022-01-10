package fileservice

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	"github.com/adamcolton/luce/tools/server/service"
)

type Handler func(path string, file *os.File, req *service.Request) *service.Response

func Wrap(fn func([]byte, *service.Request) *service.Response) Handler {
	return func(path string, f *os.File, req *service.Request) *service.Response {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return req.ResponseString(err.Error())
		}
		return fn(b, req)
	}
}

var (
	HandleBinary = Wrap(func(b []byte, req *service.Request) *service.Response {
		return req.Response(b)
	})
)

func HandleDir(path string, f *os.File, req *service.Request) *service.Response {
	dc, err := GetDirContents(f)
	if err != nil {
		return req.ResponseError(err, 500)
	}
	resp := req.Response(nil)
	err = json.NewEncoder(resp).Encode(dc)
	if err != nil {
		return req.ResponseError(err, 500)
	}
	return resp
}

type DirContents struct {
	SubDirs []string
	Files   []string
}

func GetDirContents(f *os.File) (*DirContents, error) {
	fs, err := f.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	fIdx := len(fs)
	out := make([]string, fIdx)
	dIdx := 0
	fIdx--
	for _, f := range fs {
		if f.IsDir() {
			out[dIdx] = f.Name()
			dIdx++
		} else {
			out[fIdx] = f.Name()
			fIdx--
		}
	}
	sort.Strings(out[:dIdx])
	sort.Strings(out[dIdx:])
	dc := &DirContents{
		SubDirs: out[:dIdx],
		Files:   out[dIdx:],
	}
	return dc, nil
}
