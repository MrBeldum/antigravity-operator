package templates

import "embed"

// FS contém todos os templates embutidos no binário final.
// Permite que o antigravity-operator funcione de forma 100% autônoma (zero runtime dependencies).
//
//go:embed rules/* session/* mcps/*
var FS embed.FS
