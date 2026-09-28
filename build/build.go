package main

import "github.com/goyek/goyek/v3"

var buildTask = goyek.Define(goyek.Task{
	Name:  "build",
	Usage: "go build with version injected",
	Action: func(a *goyek.A) {
		ver := resolveVersion(a)
		a.Logf("Building with version %s", ver)
		ldflags := "-s -w -X main.version=" + ver
		Exec(a, dirRoot, "go", "build", "-trimpath", "-ldflags", ldflags, "-o", "bin/server", "./cmd/server")
	},
})
