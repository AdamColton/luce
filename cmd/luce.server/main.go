package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store/bstore"
	"github.com/adamcolton/luce/tools/server"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/ltmpl"
	"github.com/quasoft/memstore"
)

type Config struct {
	Session       []string
	BoltFile      string
	Socket        string
	ServiceSocket string
	Addr          string
	Templates     struct {
		Globs lfile.MultiGlob
		lfile.PathLength
	}
	TemplateNames server.TemplateNames
	Host          string
}

func (c Config) SessionBytes() [][]byte {
	var err error
	out := make([][]byte, len(c.Session))
	for i, s := range c.Session {
		if i%2 == 1 {
			out[i], err = base64.URLEncoding.DecodeString(s)
			lerr.Panic(err)
		} else {
			out[i] = []byte(s)
		}
	}
	return out
}

func main() {
	lerr.LogTo = func(err error) {
		_, file, line, _ := runtime.Caller(2)
		fmt.Printf("%s (%d): %s", file, line, err.Error())
	}

	cfgLoc := os.Getenv("luce_server_config")
	if cfgLoc == "" {
		cfgLoc = "config.json"
	}
	r, err := os.Open(cfgLoc)
	lerr.Panic(lerr.Wrap(err, "Config Location: %s", cfgLoc))
	conf := &Config{}
	err = json.NewDecoder(r).Decode(conf)
	lerr.Panic(err)

	srvConf := &server.Config{
		Addr:          conf.Addr,
		TemplateNames: conf.TemplateNames,
		Socket:        conf.Socket,
		ServiceSocket: conf.ServiceSocket,
		UserStore:     bstore.Factory(conf.BoltFile, 0777, nil),
		SessionStore:  memstore.NewMemStore(conf.SessionBytes()...),
	}

	srvConf.Templates, err = (&ltmpl.HTMLLoader{
		Trimmer:        conf.Templates.PathLength,
		IteratorSource: conf.Templates.Globs,
	}).Load()
	lerr.Log(err)

	s, err := srvConf.New()
	lerr.Panic(err)
	s.Run()
}
