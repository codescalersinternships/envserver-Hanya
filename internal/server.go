package envserver

import (
	"encoding/json"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/env", envHandler)
	mux.HandleFunc("/env/", envHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", server.port), mux)
	return err
}

// envHandler handles server routes
func envHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	key := strings.TrimPrefix(r.URL.Path, "/env")
	switch key {
	case "":
		getEnv(w, r)
	default:
		getKeyValue(w, r)
	}
}

// getEnv displays all environment variables
func getEnv(w http.ResponseWriter, _ *http.Request) {
	encoder := json.NewEncoder(w)
	envVars := os.Environ()
	w.WriteHeader(http.StatusOK)
	err := encoder.Encode(&envVars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// getKeyValue displays key's value
func getKeyValue(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/env/")
	encoder := json.NewEncoder(w)
	value := os.Getenv(key)
	if value == "" {
		w.WriteHeader(http.StatusNotFound)
		errMessage := "key not found"
		err := encoder.Encode(&errMessage)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	err := encoder.Encode(&value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
