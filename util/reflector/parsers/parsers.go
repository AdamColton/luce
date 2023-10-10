package parsers

import (
	"strconv"
)

// TODO: Add all base types

func String(out *string, in string) (err error) {
	*out = in
	return nil
}

func Float64(f *float64, s string) (err error) {
	*f, err = strconv.ParseFloat(s, 64)
	return
}

func Int(i *int, s string) (err error) {
	*i, err = strconv.Atoi(s)
	return
}

func Int64(i *int64, s string) (err error) {
	*i, err = strconv.ParseInt(s, 0, 64)
	return
}

func Bool(i *bool, s string) (err error) {
	*i = s == "y" || s == "Y"
	return
}
