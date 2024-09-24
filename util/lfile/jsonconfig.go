package lfile

import (
	"encoding/json"
	"os"
)

func JsonConfig(env, dflt string, config any) error {
	cfgLoc := dflt
	if env != "" {
		cfgLoc = os.Getenv(env)
		if cfgLoc == "" {
			cfgLoc = dflt
		}
	}
	r, err := os.Open(cfgLoc)
	if err != nil {
		return err
	}
	return json.NewDecoder(r).Decode(config)
}
