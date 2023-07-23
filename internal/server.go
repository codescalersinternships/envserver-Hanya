package envserver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Server struct contains the port number
type Server struct {
	port int
}

// ErrInvalidPort is raised when a reserved or out-of-range port number is used
var ErrInvalidPortNum = errors.New("port number should be between 1024 and 65535")

// NewServer creates a new Server struct
func NewServer(port int) (*Server, error) {
	if !(port >= 1024 && port <= 65535) {
		return nil, ErrInvalidPortNum
	}
	server := Server{port: port}
	return &server, nil
}

// Run sets the handler for each route and runs server
func (server *Server) Run() error {
	http.HandleFunc("/env", envHandler)
	http.HandleFunc("/env/", envKeyHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil)
	return err
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	for _, e := range os.Environ() {
		fmt.Fprintln(w, e)
	}
}

func envKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/env/")
	value := os.Getenv(key)
	if value == "" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, value)
}
