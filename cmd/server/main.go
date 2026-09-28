// Command server is the HTTP service used as the subject
// of the build & deploy automation case study.
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/rifkiandriyanto/go-build-deploy-pipeline/internal/api"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	srv := &http.Server{
		Addr:              *addr,
		Handler:           api.New(version),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("my-project %s listening on %s", version, *addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
