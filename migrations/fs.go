// Package migrations will do schema migrations
package migrations

import "embed"

// go:embed *.sql

var FS embed.FS
