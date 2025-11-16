package migrations

import "embed"

// Files exposes the embedded SQL migration scripts.
//
//go:embed *.sql
var Files embed.FS
