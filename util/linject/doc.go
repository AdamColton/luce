// Package linject allows logic to be injected based on type. Note for brevity
// the term function can be a func or a method.
//
// While it would be possible to create a more generalized injector, linject is
// built around the assumption that the last argument to a function is a pointer
// to a struct. This is the "target".
//
// An injector works on a pointer to a struct. It can perform logic before a
// function or method is called, such as setting a field in the struct. Or it
// can execute logic after the call by returning a callback. There are helpers
// for targeting a single field, or custom injectors can be written for more
// complex situations.
//
// An Injector is created from an Initilizer. The Initilizer is given a
// reflect.Type. This will always be a pointer to a struct. If the struct
// fulfills the criteria the Initilizer is looking for, it will initilize an
// Injector.
//
// A collection of injection logic is represented as a slice of Initilizers. For
// a given reflect.Type they can create a slice of Injectors. And for function
// or method whose last argument is a pointer to a struct, it can create an
// injector for that function. This means that given a set of
// FunctionInitilizers fi and a function:
//  func Foo(a A, b B, s *struct{
//      Field1 T1
//      Field2 T1
//  } (string, bool)
// we can apply the initilizers as
//  injectFoo := fi.Apply(Foo).Interface().(func(A,B)(string, bool))
//
// To create an injector it generally going to be easier to call
// NewFieldInjector than to fulfill Initilizer directly.

package linject
