package main

import "html/template"

//Generated file - DO NOT EDIT

var bar = template.Must(template.New("testTmpl").Parse(`
{{define "test.txt" -}}
this is a test
{{- end}}`))
