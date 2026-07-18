// internal/share/server.go
package share

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
)

// Server is a running file-share HTTP server.
type Server struct {
	httpServer *http.Server
	addr       string
}

// Start creates the share directory if missing, then starts an HTTP
// file server for it on the given port, in a background goroutine (this
// call returns immediately - it does NOT block, unlike the old Python
// reference's serve_forever(), so callers can manage the server's
// lifecycle explicitly via Stop()). If password is non-empty, every
// request must present it via HTTP Basic Auth (username is ignored).
func Start(dir string, port int, password string) (*Server, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	fileHandler := http.FileServer(http.Dir(dir))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !checkAuth(r.Header.Get("Authorization"), password) {
			w.Header().Set("WWW-Authenticate", basicAuthChallenge())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fileHandler.ServeHTTP(w, r)
	})

	addr := portToAddr(port)
	httpServer := &http.Server{Addr: addr, Handler: mux}

	ln, err := listen(addr)
	if err != nil {
		return nil, err
	}

	go func() {
		_ = httpServer.Serve(ln) // returns http.ErrServerClosed on graceful Stop()
	}()

	return &Server{httpServer: httpServer, addr: addr}, nil
}

// Stop gracefully shuts down the server, waiting for in-flight requests
// to finish (bounded by ctx).
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Addr returns the address (host:port) the server is bound to.
func (s *Server) Addr() string {
	return s.addr
}

func portToAddr(port int) string {
	return fmt.Sprintf("0.0.0.0:%d", port)
}

func listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
