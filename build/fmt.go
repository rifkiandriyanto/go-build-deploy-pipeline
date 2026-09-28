package main

import "github.com/goyek/goyek/v3"

var fmtTask = goyek.Define(goyek.Task{
	Name:  "fmt",
	Usage: "gofmt",
	Action: func(a *goyek.A) {
		Exec(a, dirRoot, "gofmt", "-w", ".")
		Exec(a, dirBuild, "gofmt", "-w", ".")
	},
})
