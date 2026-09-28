# go-build-deploy-pipeline

A build & deploy pipeline written with Go.

[![CI](https://github.com/rifkiandriyanto/go-build-deploy-pipeline/actions/workflows/ci.yml/badge.svg)](https://github.com/rifkiandriyanto/go-build-deploy-pipeline/actions)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Case Study: Build & Deploy Pipeline with goyek

This repository is a case study demonstrating how to build a **build & deploy
pipeline** for a Go HTTP service using [goyek](https://github.com/goyek/goyek) —
a task automation library that expresses pipelines as ordinary Go code instead
of Makefiles, DSLs, or YAML.

The repository contains:

1. **Application**: an HTTP service (`cmd/server`) with `healthz`, `version`,
   and `hello` endpoints, plus a tested `internal/api` package.
2. **Automation**: a pipeline in `build/` covering the full lifecycle:
   formatting, linting, testing, binary build, Docker image build, pushing the
   image to a registry, and git-tag based releases.

---

## Problem Being Solved

Teams building HTTP services in Go usually face repetitive work that must stay
consistent across every developer machine and in CI:

- Running `gofmt`, `go vet`, and `go test` with the same flags.
- Injecting the **version** into the binary (via `-ldflags`) so that `/version`
  reports the correct value.
- Building a **Docker image** with the same version.
- Pushing the image to a **registry** for deployment.
- Creating a **release** (git tag) without missing any step.

The classic problem: these steps are scattered across Makefiles, shell scripts,
and CI configurations that differ from each other — easy to get wrong, hard to
test, and not portable across platforms.

**The goyek solution**: write the pipeline as ordinary Go functions. One
language, one way of thinking, debuggable with the Go debugger, and running the
same on Windows, macOS, and Linux.

---

## Design Decisions

### 1. Pipeline as Go code

Every task is a `goyek.Task` with `Name`, `Usage`, `Action`, and `Deps`.
Dependencies form a graph that is executed automatically in the correct order
with deduplication (a task referenced twice runs only once).

```go
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
```

### 2. Task graph & pipeline flow

```
all = mod -> fmt -> lint -> test -> build
ci  = all + diff
release = ci + push -> image
clean (standalone)
```

- `all` — local pipeline for developers.
- `ci` — CI pipeline; adds `diff`, which verifies the working tree is clean
  after all tasks run (detecting unformatted/generated files).
- `image` — builds the Docker image; **skips automatically** when the Docker
  daemon is unavailable.
- `push` — pushes the image to the registry; depends on `image`.
- `release` — marks a release: requires `-version`, verifies `ci`, then `push`.

### 3. Consistent version injection

The version comes from a single source and is used for both the binary and the
Docker image:

```go
func resolveVersion(a *goyek.A) string {
	if *version != "" { return strings.TrimPrefix(*version, "v") }
	if out, err := exec.CommandContext(a.Context(), "git", "describe", ...); err == nil {
		return strings.TrimSpace(string(out))
	}
	return "dev"
}
```

- Without flags: uses the latest git tag (e.g. `1.2.3`).
- On release: `-version=v1.2.3` guarantees an explicit version.
- Binary: `-ldflags "-X main.version=..."` → returned by `/version`.
- Docker: `--build-arg VERSION=...` → injected in the `Dockerfile`.

### 4. Graceful degradation across environments

Tasks that need Docker check for the daemon and call `a.Skip(...)` instead of
`a.Fatal(...)`. As a result, the `image` pipeline can still be run and tested on
machines without Docker (e.g. a laptop), while in CI with a Docker daemon the
task runs fully.

### 5. Middleware for a good developer experience

```go
goyek.UseExecutor(middleware.ReportFlow)
goyek.Use(middleware.ReportStatus)
goyek.Use(middleware.ReportLongRun(time.Minute))
if !*v { goyek.Use(middleware.SilentNonFailed) }
if *dryRun { goyek.Use(middleware.DryRun) }
```

The result: concise output (only failing tasks are shown), colors, per-task
status, and a `-dry-run` mode to preview without executing anything.

---

## Directory Layout

```
go-build-deploy-pipeline/
├── cmd/server/          # HTTP service (application entrypoint)
├── internal/api/        # HTTP handlers + unit tests
├── build/               # AUTOMATION MODULE (own go.mod)
│   ├── main.go          # CLI entry, flags, middleware
│   ├── exec.go          # helper for running external commands
│   ├── version.go       # resolveVersion + docker availability check
│   ├── all.go / ci.go   # aggregate pipelines
│   ├── mod.go / fmt.go / lint.go / test.go / build.go
│   ├── image.go / push.go / release.go
│   └── clean.go / diff.go
├── Dockerfile           # multi-stage build (distroless)
├── goyek.sh             # bash wrapper
├── goyek.ps1            # PowerShell wrapper
└── .github/workflows/ci.yml
```

### Why does `build/` have its own `go.mod`?

The automation module is separated from the application module so that tooling
dependencies (here, just goyek) do not pollute the application's `go.mod`. This
is the pattern goyek recommends and makes the pipeline easy to move to other
repositories.

---

## Usage

```sh
./goyek.sh all            # full local pipeline
./goyek.sh test -v        # individual task, verbose
./goyek.sh ci             # CI pipeline (includes working tree check)
./goyek.sh all -dry-run   # preview without executing
./goyek.sh -h             # list all tasks & flags
```

Example: run the built binary:

```sh
./goyek.sh build
./bin/server &
curl http://localhost:8080/healthz   # ok
curl http://localhost:8080/version   # {"version":"1.2.3"}
curl "http://localhost:8080/hello?name=goyek"  # {"message":"Hello, goyek!"}
```

### Release

```sh
./goyek.sh release -version=v1.2.3
```

This runs `ci` (full verification), builds the image with version `1.2.3`,
pushes it to `$REGISTRY`, and creates the git tag `v1.2.3`.

### Pipeline output

```out
===== TASK  test
      test.go:14: Run [go test -race -covermode=atomic -coverprofile=coverage.out ./...] in .
ok  	github.com/rifkiandriyanto/go-build-deploy-pipeline/internal/api	0.574s
----- PASS: test (1.49s)
===== TASK  build
      build.go:9: Run [go build -trimpath -ldflags "-s -w -X main.version=..." -o bin/server ./cmd/server] in .
----- PASS: build (0.06s)
```

---

## Comparison with Alternatives

| Aspect | goyek | Make | Mage |
|---|---|---|---|
| Language | Go | Shell | Go |
| IDE debugging | Yes | No | Yes |
| Cross-platform | Yes | Partial | Yes |
| Dependency graph | Yes (Deps) | Yes | Yes |
| Dry-run mode | Yes (middleware) | No | No |
| Tool installation | Not needed | Not needed | Binary needed |

---

## Takeaways

1. **A single source of truth for the version** — injecting the version into
   the binary and the image from one function removes the whole class of "image
   does not match binary" bugs.
2. **A testable pipeline** — since the pipeline is Go, helpers such as
   `dockerAvailable` and `resolveVersion` can be tested directly.
3. **Graceful degradation** — `a.Skip` lets the same pipeline run on a
   developer laptop and in CI, across different environments.
4. **Output ownership** — every command writes through `a.Output()`, so
   middleware can handle buffering/parallelism without races.
5. **Composability** — small tasks (`fmt`, `lint`, `test`) compose into
   pipelines (`all`, `ci`, `release`) through `Deps`, so they can be run
   individually or together.

---

## License

[MIT](LICENSE) — Copyright (c) 2026 [Rifki Andriyanto](https://github.com/rifkiandriyanto)