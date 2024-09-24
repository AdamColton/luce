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

type ConfigHandlers struct {
	Usage struct {
		Show   string
		Update string
	}
	Disabled bool
	reflector.Parser[string]
	config any
}

func NewConfigHandlers(config any) (*ConfigHandlers, error) {
	if !ltype.IsPtrToStruct.OnInterface(config) {
		return nil, lerr.Str("must be pointer to struct")
	}
	ch := &ConfigHandlers{
		config: config,
		Parser: Parser,
	}
	ch.Usage.Show = "show current config values"
	return ch, nil
}

type ShowConfigReq struct{}

type ShowConfigResp struct {
	Config slice.Slice[[2]string]
}

func (ch ConfigHandlers) ShowConfigHandler(req *ShowConfigReq) *ShowConfigResp {
	v := reflect.ValueOf(ch.config).Elem()
	t := v.Type()
	fs := v.NumField()
	out := make(slice.Slice[[2]string], fs)
	for i := 0; i < fs; i++ {
		out[i][0] = t.Field(i).Name
		out[i][1] = fmt.Sprint(v.Field(i).Interface())
	}

	out.Sort(func(i, j [2]string) bool {
		return i[0] < j[0]
	})

	return &ShowConfigResp{Config: out}
}

func (ch ConfigHandlers) ShowConfigUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    ch.Usage.Show,
		Disabled: ch.Disabled,
	}
}

func (rnr *Runner) ShowConfigRespHandler(cfg *ShowConfigResp) {
	out := make([]string, len(cfg.Config))
	for i, f := range cfg.Config {
		out[i] = fmt.Sprintf("%s: %s", f[0], f[1])
	}
	rnr.WriteString(strings.Join(out, "\n"))
}

type UpdateConfigReq struct {
	Field, Value string
}

type UpdateConfigResp struct {
	Msg string
}

func (ch ConfigHandlers) UpdateConfigHandler(req *UpdateConfigReq) *UpdateConfigResp {
	v := reflect.ValueOf(ch.config).Elem()
	f := v.FieldByName(req.Field)
	if f.Kind() == reflect.Invalid {
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
		Msg: fmt.Sprintf("updated"),
	}
}

func (ch ConfigHandlers) UpdateConfigUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    ch.Usage.Update,
		Disabled: ch.Disabled,
	}
}

func (rnr *Runner) UpdateConfigRespHandler(u *UpdateConfigResp) {
	rnr.WriteString(u.Msg)
}
