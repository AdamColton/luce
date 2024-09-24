package lfile

import (
	"encoding/json"
	"os"
)

// JsonConfig decodes the JSON file that holds the configuration into config,
// which should be a pointer. The file is read from the operating system. If env
// is not empty and the environment variable it names is set, the value of that
// variable is the path of the file, otherwise the path is dflt. It is meant to be
// called with a service's own variable and a default file name, for example
// JsonConfig("my_service_config", "my_service.json", &cfg). The file is closed
// before it returns.
func JsonConfig(env, dflt string, config any) error {
	return JsonConfigFS(OSRepository{}, env, dflt, config)
}

// JsonConfigFS is JsonConfig reading the file from o. The environment variable
// still comes from the operating system; its value is a path in o.
func JsonConfigFS(o FSOpener, env, dflt string, config any) error {
	cfgLoc := dflt
	if env != "" {
		cfgLoc = os.Getenv(env)
		if cfgLoc == "" {
			cfgLoc = dflt
		}
	}
	f, err := o.Open(cfgLoc)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(config)
}
