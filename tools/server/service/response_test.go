package service_test

import (
	"html/template"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}
	body := []byte("test body")
	resp := req.Response(body).
		SetStatus(205)
	assert.Equal(t, body, resp.Body)
	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, service.ResponseTypeID32, resp.TypeID32())
	assert.Equal(t, 205, resp.Status)
}

func TestResponseString(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}
	body := "test body"
	resp := req.ResponseString(body)
	assert.Equal(t, []byte(body), resp.Body)
	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, service.ResponseTypeID32, resp.TypeID32())
}

func TestResponseErr(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}
	body := "test error"
	err := lerr.Str(body)
	resp := req.ResponseErr(err, 500)
	assert.Equal(t, []byte(body), resp.Body)
	assert.Equal(t, req.ID, resp.ID)
}

type statusErrWrapper struct {
	err    error
	status int
}

func (ser statusErrWrapper) Error() string {
	return ser.err.Error()
}

func (ser statusErrWrapper) Status() int {
	return ser.status
}

func TestErrCheck(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}
	body := "test error"
	var err error = lerr.Str(body)
	resp := req.ErrCheck(err)
	assert.Equal(t, []byte(body), resp.Body)
	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, 500, resp.Status)

	err = statusErrWrapper{
		err:    err,
		status: 501,
	}
	resp = req.ErrCheck(err)
	assert.Equal(t, []byte(body), resp.Body)
	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, 501, resp.Status)

	resp = req.ErrCheck(nil)
	assert.Nil(t, resp)
}

type Person struct {
	Name string
	Age  int
}

func TestSerializeResponse(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}

	person := Person{
		Name: "Adam",
		Age:  40,
	}
	s := json.Serializer{}

	resp, err := req.SerializeResponse(s, person, nil)
	assert.NoError(t, err)

	expected, err := s.Serialize(person, nil)
	assert.NoError(t, err)

	assert.Equal(t, expected, resp.Body)
}

func TestResponseTemplate(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}

	person := Person{
		Name: "Adam",
		Age:  40,
	}

	tmpl := template.New("root")
	_, err := tmpl.Parse(`Name: {{.Name}} Age: {{.Age}}{{define "named"}}Named.Name: {{.Name}} Named.Age: {{.Age}}{{end}}`)
	assert.NoError(t, err)

	resp := req.ResponseTemplate("", tmpl, person)
	expected := `Name: Adam Age: 40`
	assert.Equal(t, []byte(expected), resp.Body)

	resp = req.ResponseTemplate("named", tmpl, person)
	expected = `Named.Name: Adam Named.Age: 40`
	assert.Equal(t, []byte(expected), resp.Body)
}

func TestResponseRedirect(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}
	url := "test-url"
	resp := req.Redirect(url)
	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, service.HttpRedirect, resp.Status)
	assert.Equal(t, []byte(url), resp.Body)
}
