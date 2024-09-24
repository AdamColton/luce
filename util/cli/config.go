package cli

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

// ConfigHandlers has the handlers for two commands, showConfig and
// updateConfig, that show and change the exported fields of a config struct. It
// is meant to be embedded in the struct that a handler.MethodsRegistrar reads.
type ConfigHandlers struct {
	// Usage is the description of each command.
	Usage struct {
		Show   string
		Update string
	}
	// Disabled leaves both commands out.
	Disabled bool
	// Parser turns the text of a new value into the type of its field.
	reflector.Parser[string]
	config any
}

// NewConfigHandlers creates ConfigHandlers for config, which must be a pointer
// to a struct. It uses Parser and sets the usage of both commands.
func NewConfigHandlers(config any) (*ConfigHandlers, error) {
	if !ltype.IsPtrToStruct.OnInterface(config) {
		return nil, lerr.Str("must be pointer to struct")
	}
	ch := &ConfigHandlers{
		config: config,
		Parser: Parser,
	}
	ch.Usage.Show = "show current config values"
	ch.Usage.Update = "update a config value"
	return ch, nil
}

// ShowConfigReq is the request for the showConfig command.
type ShowConfigReq struct{}

// ShowConfigResp is the response to the showConfig command.
type ShowConfigResp struct {
	// Config is the name and value of each exported field, ordered by name.
	Config slice.Slice[[2]string]
}

// ShowConfigHandler answers the showConfig command with the exported fields of
// the config.
func (ch ConfigHandlers) ShowConfigHandler(req *ShowConfigReq) *ShowConfigResp {
	v := reflect.ValueOf(ch.config).Elem()
	t := v.Type()
	fs := v.NumField()
	out := make(slice.Slice[[2]string], 0, fs)
	for i := 0; i < fs; i++ {
		if t.Field(i).IsExported() {
			out = append(out, [2]string{t.Field(i).Name, fmt.Sprint(v.Field(i).Interface())})
		}
	}

	out.Sort(func(i, j [2]string) bool {
		return i[0] < j[0]
	})

	return &ShowConfigResp{Config: out}
}

// ShowConfigUsage describes the showConfig command.
func (ch ConfigHandlers) ShowConfigUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    ch.Usage.Show,
		Disabled: ch.Disabled,
	}
}

// ShowConfigRespHandler is the handler for a ShowConfigResp. It writes a line
// for each field, as "name: value".
func (rnr *Runner) ShowConfigRespHandler(cfg *ShowConfigResp) {
	out := make([]string, len(cfg.Config))
	for i, f := range cfg.Config {
		out[i] = fmt.Sprintf("%s: %s", f[0], f[1])
	}
	rnr.WriteString(strings.Join(out, "\n"))
}

// UpdateConfigReq is the request for the updateConfig command. Value is parsed
// into the field named Field.
type UpdateConfigReq struct {
	Field, Value string
}

// UpdateConfigResp is the response to the updateConfig command. Msg says what
// happened.
type UpdateConfigResp struct {
	Msg string
}

// UpdateConfigHandler answers the updateConfig command. The message is
// "updated", the error from parsing the value, or that the field is not defined.
// A field that is not exported is not defined.
func (ch ConfigHandlers) UpdateConfigHandler(req *UpdateConfigReq) *UpdateConfigResp {
	v := reflect.ValueOf(ch.config).Elem()
	f := v.FieldByName(req.Field)
	if f.Kind() == reflect.Invalid || !f.CanSet() {
		return &UpdateConfigResp{
			Msg: fmt.Sprintf("field:%s is not defined", req.Field),
		}
	}

	err := ch.Parse(f.Addr().Interface(), req.Value)
	if err != nil {
		return &UpdateConfigResp{
			Msg: err.Error(),
		}
	}
	return &UpdateConfigResp{
		Msg: "updated",
	}
}

// UpdateConfigUsage describes the updateConfig command.
func (ch ConfigHandlers) UpdateConfigUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    ch.Usage.Update,
		Disabled: ch.Disabled,
	}
}

// UpdateConfigRespHandler is the handler for an UpdateConfigResp. It writes the
// message.
func (rnr *Runner) UpdateConfigRespHandler(u *UpdateConfigResp) {
	rnr.WriteString(u.Msg)
}
