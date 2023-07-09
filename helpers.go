package govictoria

import "encoding/base64"

// BasicAuth returns the base64 encoded string for basic auth
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
