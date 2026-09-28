package main

import "github.com/goyek/goyek/v3"

var ci = goyek.Define(goyek.Task{
	Name:  "ci",
	Usage: "CI pipeline (all + diff)",
	Deps: goyek.Deps{
		all,
		diff,
	},
})
