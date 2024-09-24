package cli_test

import (
	"testing"

	"github.com/adamcolton/luce/util/cli"
	"github.com/stretchr/testify/assert"
)

type appConfig struct {
	Name   string
	Volume int
	On     bool
	secret string
}

func TestNewConfigHandlersErrors(t *testing.T) {
	for name, config := range map[string]any{
		"not a pointer":      appConfig{},
		"not a struct":       new(int),
		"nil":                nil,
		"pointer to pointer": new(*appConfig),
	} {
		t.Run(name, func(t *testing.T) {
			ch, err := cli.NewConfigHandlers(config)
			assert.Nil(t, ch)
			assert.EqualError(t, err, "must be pointer to struct")
		})
	}
}

func TestShowConfig(t *testing.T) {
	cfg := &appConfig{Name: "app", Volume: 11, secret: "hidden"}
	ch, err := cli.NewConfigHandlers(cfg)
	assert.NoError(t, err)

	// The exported fields, in order.
	resp := ch.ShowConfigHandler(&cli.ShowConfigReq{})
	assert.Equal(t, [][2]string{{"Name", "app"}, {"On", "false"}, {"Volume", "11"}}, [][2]string(resp.Config))

	details := ch.ShowConfigUsage()
	assert.Equal(t, "show current config values", details.Usage)
	assert.False(t, details.Disabled)
	ch.Disabled = true
	assert.True(t, ch.ShowConfigUsage().Disabled)
	assert.True(t, ch.UpdateConfigUsage().Disabled)
}

func TestUpdateConfig(t *testing.T) {
	cfg := &appConfig{Name: "app"}
	ch, err := cli.NewConfigHandlers(cfg)
	assert.NoError(t, err)
	assert.Equal(t, "update a config value", ch.UpdateConfigUsage().Usage)

	update := func(field, value string) string {
		return ch.UpdateConfigHandler(&cli.UpdateConfigReq{Field: field, Value: value}).Msg
	}
	assert.Equal(t, "updated", update("Volume", "7"))
	assert.Equal(t, "updated", update("Name", "renamed"))
	assert.Equal(t, "updated", update("On", "y"))
	assert.Equal(t, &appConfig{Name: "renamed", Volume: 7, On: true}, cfg)

	// A value that can't be parsed is reported.
	assert.Contains(t, update("Volume", "loud"), "invalid syntax")

	// Fields that are not exported, or are not there, are not defined.
	assert.Equal(t, "field:secret is not defined", update("secret", "x"))
	assert.Equal(t, "field:Nope is not defined", update("Nope", "x"))
	assert.Equal(t, "", cfg.secret)
}

func TestConfigResponses(t *testing.T) {
	r, buf, _ := newRunner(t, nil)

	r.ShowConfigRespHandler(&cli.ShowConfigResp{Config: [][2]string{{"Name", "app"}, {"Volume", "11"}}})
	assert.Equal(t, "Name: app\nVolume: 11", buf.String())

	buf.Reset()
	r.UpdateConfigRespHandler(&cli.UpdateConfigResp{Msg: "updated"})
	assert.Equal(t, "updated", buf.String())
}
