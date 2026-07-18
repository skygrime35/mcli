package share

import (
	"encoding/base64"
	"testing"
)

func basicAuthHeader(user, pass string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+pass))
}

func TestCheckAuth_NoPasswordConfigured(t *testing.T) {
	// An empty configured password means auth is disabled entirely -
	// any request (even with no Authorization header) is allowed.
	if !checkAuth("", "") {
		t.Error("expected checkAuth to allow access when no password is configured")
	}
	if !checkAuth(basicAuthHeader("anyone", "wrong"), "") {
		t.Error("expected checkAuth to allow access when no password is configured, regardless of header")
	}
}

func TestCheckAuth_CorrectPassword(t *testing.T) {
	if !checkAuth(basicAuthHeader("ignored-username", "secret123"), "secret123") {
		t.Error("expected checkAuth to succeed with the correct password")
	}
}

func TestCheckAuth_UsernameIsIgnored(t *testing.T) {
	if !checkAuth(basicAuthHeader("completely-different-user", "secret123"), "secret123") {
		t.Error("expected checkAuth to succeed regardless of username, matching the old reference's design")
	}
}

func TestCheckAuth_WrongPassword(t *testing.T) {
	if checkAuth(basicAuthHeader("user", "wrong"), "secret123") {
		t.Error("expected checkAuth to reject an incorrect password")
	}
}

func TestCheckAuth_MissingHeader(t *testing.T) {
	if checkAuth("", "secret123") {
		t.Error("expected checkAuth to reject a missing Authorization header when a password IS configured")
	}
}

func TestCheckAuth_MalformedHeader(t *testing.T) {
	if checkAuth("Bearer sometoken", "secret123") {
		t.Error("expected checkAuth to reject a non-Basic Authorization header")
	}
	if checkAuth("Basic not-valid-base64!!!", "secret123") {
		t.Error("expected checkAuth to reject unparseable base64")
	}
}
