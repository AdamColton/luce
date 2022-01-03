#!/bin/bash

rm base.go
printf "package filter\n\n" >> base.go
luce gen f int >> base.go
luce gen f string >> base.go
luce gen f float64 -t Float >> base.go