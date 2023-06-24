// luce.server runs an instance of luce/tools/server. It uses quasoft/memstore
// for the SessionStore. It uses a boltdb/bolt database for the UserStore.
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store/bstore"
	"github.com/adamcolton/luce/tools/server"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/ltmpl"
	"github.com/quasoft/memstore"
)

// Config holds the values needed to create a Luce Server. If an environmental
// variable is set for luce_server_config, it will look in that location for a
// config file. Otherwise it will look in the running directory for config.json.
// The config file should be json formatted.
type Config struct {
	// Session is a list of base 64 URL Encoded key pairs.
	Session []string
	// BoltFile is the file used for the bolt database
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

// SessionBytes converts Session into a format that memstore.NewMemStore can
// use.
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

	ss := memstore.NewMemStore(conf.SessionBytes()...)
	ss.Options.Domain = conf.Host

	srvConf := &server.Config{
		Addr:          conf.Addr,
		TemplateNames: conf.TemplateNames,
		Socket:        conf.Socket,
		ServiceSocket: conf.ServiceSocket,
		UserStore:     bstore.Factory(conf.BoltFile, 0777, nil),
		SessionStore:  ss,
		Host:          conf.Host,
	}

	srvConf.Templates, err = (&ltmpl.HTMLLoader{
		Trimmer:        conf.Templates.PathLength,
		IteratorSource: conf.Templates.Globs,
	}).Load()
	lerr.Log(err)

	s, err := srvConf.New()
	lerr.Panic(err)

	go func() {
		rdr := iobus.Config{
			Sleep: time.Millisecond,
		}.NewReader(os.Stdin)
		ctx := cli.NewContext(os.Stdout, rdr.Out, nil)
		s.Cli(ctx, nil)
	}()

	s.Run()
}
