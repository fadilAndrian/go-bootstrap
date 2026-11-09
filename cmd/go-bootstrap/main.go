package main

import (
	"github.com/konnco/go-bootstrap/internal/cli"

	// Statically register built-in plugins
	_ "github.com/konnco/go-bootstrap/plugins/database/mysql"
	_ "github.com/konnco/go-bootstrap/plugins/database/postgres"
	_ "github.com/konnco/go-bootstrap/plugins/database/sqlite"
	_ "github.com/konnco/go-bootstrap/plugins/http/chi"
	_ "github.com/konnco/go-bootstrap/plugins/http/gin"
	_ "github.com/konnco/go-bootstrap/plugins/http/nethttp"
	_ "github.com/konnco/go-bootstrap/plugins/worker/none"
	_ "github.com/konnco/go-bootstrap/plugins/worker/river"
)

func main() {
	cli.Execute()
}
