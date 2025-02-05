package https

import (
	"embed"
)

//go:embed certs
var certs embed.FS

// GetTLSCerts GetTLSCerts
func GetTLSCerts() embed.FS {
	return certs
}
