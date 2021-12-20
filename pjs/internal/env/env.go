package env

import (
	"os"
	"strings"
)

// ParseAlternateSettings searches for alternative postgres keys.
// Not otherwise checked by PGX's libraries.
func ParseAlternateSettings() string {
	nameMap := map[string]string{
		"POSTGRES_HOST":     "host",
		"POSTGRES_PORT":     "port",
		"POSTGRES_DB":       "database",
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_SSLMODE":  "sslmode",
		"POSTGRES_ROOTCERT": "sslrootcert",
	}

	var builder strings.Builder
	for envname, realname := range nameMap {
		value := os.Getenv(envname)
		if value != "" {
			builder.WriteString(realname + "=" + value + " ")
		}
	}
	return builder.String()
}
