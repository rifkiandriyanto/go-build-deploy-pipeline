package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/goyek/goyek/v3"
)

// resolveVersion returns the version to use for build and release.
// It prefers the -version flag, then the latest git tag, then "dev".
func resolveVersion(a *goyek.A) string {
	a.Helper()
	if *version != "" {
		return strings.TrimPrefix(*version, "v")
	}
	if out, err := exec.CommandContext(a.Context(), "git", "describe", "--tags", "--always", "--dirty").Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	return "dev"
}

// dockerAvailable reports whether a working Docker CLI is present.
func dockerAvailable(a *goyek.A) bool {
	a.Helper()
	_, err := exec.LookPath("docker")
	if err != nil {
		a.Log("docker not found")
		return false
	}
	cmd := exec.CommandContext(a.Context(), "docker", "info")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}
