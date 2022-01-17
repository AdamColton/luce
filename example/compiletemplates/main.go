package main

import (
	"bytes"
	"fmt"
)

//go:generate compiletemplates

func main() {
	buf := bytes.NewBuffer(nil)
	bar.ExecuteTemplate(buf, "test.txt", nil)
	fmt.Println(buf.String())
}
