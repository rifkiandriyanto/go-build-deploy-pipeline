package main

import "github.com/goyek/goyek/v3"

// release creates a git tag and publishes the docker image.
// Usage: goyek.sh release -version=v1.2.3
var release = goyek.Define(goyek.Task{
	Name:  "release",
	Usage: "tag a release and publish image (requires -version)",
	Deps: goyek.Deps{
		ci,
		push,
	},
	Action: func(a *goyek.A) {
		ver := resolveVersion(a)
		if *version == "" {
			a.Fatal("release requires -version flag, e.g. -version=v1.2.3")
		}
		a.Logf("Tagging release v%s", ver)
		Exec(a, dirRoot, "git", "tag", "-a", "v"+ver, "-m", "Release v"+ver)
		Exec(a, dirRoot, "git", "push", "origin", "v"+ver)
	},
})
