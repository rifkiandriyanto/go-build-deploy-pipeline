package main

import "github.com/goyek/goyek/v3"

var image = goyek.Define(goyek.Task{
	Name:  "image",
	Usage: "build docker image (requires docker daemon)",
	Action: func(a *goyek.A) {
		if !dockerAvailable(a) {
			a.Skip("docker is not available")
			return
		}
		ver := resolveVersion(a)
		a.Logf("Building image %s:%s", *registry, ver)
		Exec(a, dirRoot, "docker", "build", "--build-arg", "VERSION="+ver, "-t", *registry+":"+ver, ".")
	},
})
