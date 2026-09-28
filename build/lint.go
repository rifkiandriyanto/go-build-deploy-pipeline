package main

import "github.com/goyek/goyek/v3"

var lint = goyek.Define(goyek.Task{
	Name:  "lint",
	Usage: "go vet",
	Action: func(a *goyek.A) {
		Exec(a, dirRoot, "go", "vet", "./...")
	},
})
