package service_test

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/stretchr/testify/assert"
)

func TestRouteChainableHelpers(t *testing.T) {
	r := service.NewRoute("foo").
		Get().
		WithQuery().
		WithUser().
		WithBody().
		WithPrefix().
		RequireGroup("admin")

	assert.Equal(t, "GET", r.Method)
	assert.True(t, r.Query)
	assert.True(t, r.User)
	assert.True(t, r.Body)
	assert.True(t, r.PathPrefix)
	assert.Equal(t, "admin", r.Require.Group)
	assert.Equal(t, []string{"GET"}, r.Methods())
	assert.Equal(t, "(GET) /foo...", r.String())

	r.RequireGroup("media")
	assert.Equal(t, "admin,media", r.Require.Group)

	r.Post()
	assert.Equal(t, "GET,POST", r.Method)

	r2 := service.NewRoute("bar").Post()
	assert.Equal(t, "POST", r2.Method)

	r3 := service.NewRoute("bar").Delete()
	assert.Equal(t, "DELETE", r3.Method)

	r4 := service.NewRoute("bar").Put()
	assert.Equal(t, "PUT", r4.Method)

	r5 := service.NewRoute("bar").WithForm()
	assert.True(t, r5.Form)

	r6 := service.NewRoute("bar").PostForm()
	assert.True(t, r6.Form)
	assert.Equal(t, "POST", r6.Method)

	err := (&service.Route{}).Validate()
	assert.Equal(t, service.ErrPathRequired, err)

	withVars := service.NewRoute("/foo/{id}")
	assert.NoError(t, withVars.Validate())
	assert.True(t, withVars.PathVars)
	assert.Equal(t, "GET", withVars.Method)
	assert.Equal(t, withVars.String(), withVars.ID)
}
