package main

import "github.com/goyek/goyek/v3"

var all = goyek.Define(goyek.Task{
	Name:  "all",
	Usage: "build pipeline (mod, fmt, lint, test, build)",
	Deps: goyek.Deps{
		mod,
		fmtTask,
		lint,
		test,
		buildTask,
	},
})
