// Package midware provides a way to convert a MidwareFunc to an
// http.HandlerFunc. MidwareFunc is not a literal type but is any func of the
// form
// - func(w http.ResponseWriter, r *http.Request, data struct{...})
//
// Midware is used to inject values into data using reflection.
package midware
