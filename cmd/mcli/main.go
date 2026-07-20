// cmd/mcli/main.go
package main

import "github.com/skygrime35/mcli/internal/cli"

// version is overridden at build time via:
//   -ldflags "-X main.version=v1.2.3"
// GoReleaser does this automatically for release builds; local `go build`
// and `go run` leave it at "dev".
var version = "dev"

func main() {
	cli.Execute(version)
}
