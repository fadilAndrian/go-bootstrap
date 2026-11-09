package main

import (
	"github.com/fadilAndrian/go-bootstrap/internal/cli"

	// Statically register built-in plugins
	_ "github.com/fadilAndrian/go-bootstrap/plugins/database/mysql"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/database/postgres"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/database/sqlite"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/http/chi"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/http/gin"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/http/nethttp"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/worker/none"
	_ "github.com/fadilAndrian/go-bootstrap/plugins/worker/river"
)

func main() {
	cli.Execute()
}
