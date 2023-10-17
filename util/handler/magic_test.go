package handler_test

import (
	"testing"

	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

// oddObj has methods that DefaultRegistrar has to be careful with.
type oddObj struct{}

// Handler has no name left once the suffix is removed.
func (oddObj) Handler(i int) {}

func (oddObj) NilUsageHandler(i int) {}
func (oddObj) NilUsageUsage() *handler.CommandDetails {
	return nil
}

func (oddObj) RenamedHandler(s string) string { return s }
func (oddObj) RenamedUsage() *handler.CommandDetails {
	return &handler.CommandDetails{Name: "other", Alias: "o", Usage: "a renamed command"}
}

// WrongUsageUsage does not return *CommandDetails, so it is not used.
func (oddObj) WrongUsageHandler(s string) {}
func (oddObj) WrongUsageUsage() string    { return "ignored" }
func (oddObj) TwoArgsHandler(a, b int)    {}
func (oddObj) NoArgsHandler()             {}

func TestCommandDetails(t *testing.T) {
	cmds := handler.DefaultRegistrar.Commands(oddObj{})

	// The names come from the methods, unless the details give one.
	assert.Equal(t, 3, cmds.Len())
	assert.Equal(t, "nilUsage", cmds.GetVal("nilUsage").Name)
	assert.Equal(t, "wrongUsage", cmds.GetVal("wrongUsage").Name)
	assert.Empty(t, cmds.GetVal("wrongUsage").Usage)

	renamed := cmds.GetVal("other")
	assert.Equal(t, "other", renamed.Name)
	assert.Equal(t, "o", renamed.Alias)
	assert.Equal(t, "a renamed command", renamed.Usage)

	// The action is the method.
	assert.Equal(t, "hello", renamed.Action.(func(string) string)("hello"))

	// A method named Handler has no name for a command, so it is disabled.
	details := handler.DefaultRegistrar.Detailer(reflector.MethodOn(oddObj{}, "Handler"))
	assert.True(t, details.Disabled)

	// It and the methods that don't take one argument are left out.
	for _, name := range []string{"", "handler", "twoArgs", "noArgs"} {
		_, found := cmds.Get(name)
		assert.False(t, found, name)
	}
}

// badObj has a method that can't be a Handler.
type badObj struct{}

func (badObj) GoodHandler(s string) string { return s }
func (badObj) BadHandler(i int) (int, int) { return i, i }

func TestRegisterErrors(t *testing.T) {
	s := handler.NewSwitch(2)
	ts, err := handler.DefaultRegistrar.Register(s, badObj{})
	assert.Error(t, err)
	assert.Len(t, ts, 1)

	// The valid method is registered anyway.
	a, err := s.Handle("still works")
	assert.NoError(t, err)
	assert.Equal(t, "still works", a)
}
