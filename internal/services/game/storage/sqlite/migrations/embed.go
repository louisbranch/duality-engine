package migrations

import "embed"

//go:embed events/*.sql
var EventsFS embed.FS

//go:embed projections/*.sql
var ProjectionsFS embed.FS

//go:embed content/*.sql
var ContentFS embed.FS
