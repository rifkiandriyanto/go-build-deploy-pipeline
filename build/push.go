package main

import "github.com/goyek/goyek/v3"

var push = goyek.Define(goyek.Task{
	Name:  "push",
	Usage: "push docker image to registry (requires docker + registry login)",
	Deps: goyek.Deps{
		image,
	},
	Action: func(a *goyek.A) {
		if !dockerAvailable(a) {
			a.Skip("docker is not available")
			return
		}
		ver := resolveVersion(a)
		a.Logf("Pushing image %s:%s", *registry, ver)
		Exec(a, dirRoot, "docker", "push", *registry+":"+ver)
	},
})
