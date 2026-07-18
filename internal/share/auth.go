package share

import (
	"crypto/subtle"
	"encoding/base64"
	"strings"
)

const basicAuthRealm = "File Share"

// checkAuth reports whether a request carrying authHeader (the raw
// "Authorization" header value, or "" if absent) may proceed, given the
// configured password. If password is empty, auth is disabled entirely
// and every request is allowed - matching this project's deliberate
// default of "no password configured = open access" (see Global
// Constraints: unlike the old Python reference, there is no hardcoded
// fallback password).
//
// Only the password portion of "user:pass" is checked, using a
// constant-time comparison; the username is intentionally ignored,
// matching the old reference's actual design (a single shared secret,
// not per-user accounts).
func checkAuth(authHeader string, password string) bool {
	if password == "" {
		return true
	}
	const prefix = "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, prefix))
	if err != nil {
		return false
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(parts[1]), []byte(password)) == 1
}

// basicAuthChallenge returns the WWW-Authenticate header value sent
// alongside a 401 response.
func basicAuthChallenge() string {
	return `Basic realm="` + basicAuthRealm + `"`
}
