package cli_test

import (
	"testing"

	"github.com/adamcolton/luce/util/cli"
	"github.com/stretchr/testify/assert"
)

func TestHTMLTemplateLoadHandler(t *testing.T) {
	calls := 0
	h := cli.NewHTMLTemplateLoadHandler(func() { calls++ })

	assert.Equal(t, &cli.HTMLTemplateLoadResp{}, h.LoadHTMLTemplateHandler(&cli.HTMLTemplateLoadReq{}))
	assert.Equal(t, 1, calls)

	// The details can be changed through the handler.
	details := h.LoadHTMLTemplateUsage()
	assert.Equal(t, "Reload HTML Templates", details.Usage)
	h.Details.Usage = "reload"
	assert.Equal(t, "reload", h.LoadHTMLTemplateUsage().Usage)
}
