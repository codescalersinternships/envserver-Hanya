package envserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	testCases := []struct {
		name         string
		port         int
		returnsError bool
	}{
		{
			name:         "valid port runs server",
			port:         3000,
			returnsError: false,
		},
		{
			name:         "invalid port raises error",
			port:         1023,
			returnsError: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			server, err := NewServer(tC.port)
			if tC.returnsError {
				assert.NotNil(t, err)
				assert.Equal(t, ErrInvalidPortNum, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.Equal(t, tC.port, server.port)
			}
		})
	}
}

func TestEnvKeyHandler(t *testing.T) {
	testCases := []struct {
		name         string
		key          string
		value        string
		returnsError bool
	}{
		{
			name:         "existing key returns correct value",
			key:          "USER",
			returnsError: false,
		},
		{
			name:         "empty value raises error 404",
			key:          "USER",
			value:        "",
			returnsError: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			value := os.Getenv(tC.key)
			if tC.returnsError {
				os.Setenv(tC.key, tC.value)
				defer os.Unsetenv(tC.key)
			}
			req, err := http.NewRequest(http.MethodGet, "/env/"+tC.key, nil)
			assert.NoError(t, err, "failed to create request")
			writer := httptest.NewRecorder()
			envHandler(writer, req)
			resp := writer.Result()
			body := strings.TrimSpace(writer.Body.String())

			body = body[1 : len(body)-1]
			if tC.returnsError {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, value, body)
			}
		})
	}
}
func TestEnvHandler(t *testing.T) {
	testCases := []struct {
		name       string
		method     string
		returnsErr bool
	}{
		{
			name:       "post returns error 403",
			method:     http.MethodPost,
			returnsErr: true,
		},
		{
			name:       "put returns error 403",
			method:     http.MethodPut,
			returnsErr: true,
		},
		{
			name:       "connect returns error 403",
			method:     http.MethodConnect,
			returnsErr: true,
		},
		{
			name:       "delete returns error 403",
			method:     http.MethodDelete,
			returnsErr: true,
		},
		{
			name:       "head returns error 403",
			method:     http.MethodHead,
			returnsErr: true,
		},
		{
			name:       "options returns error 403",
			method:     http.MethodOptions,
			returnsErr: true,
		},
		{
			name:       "patch returns error 403",
			method:     http.MethodPatch,
			returnsErr: true,
		},
		{
			name:       "trace returns error 403",
			method:     http.MethodTrace,
			returnsErr: true,
		},
		{
			name:       "get returns correct json",
			method:     http.MethodGet,
			returnsErr: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			req, err := http.NewRequest(tC.method, "/env", nil)
			assert.NoError(t, err, "failed to create request")
			res := httptest.NewRecorder()
			envHandler(res, req)
			if tC.returnsErr {
				assert.Equal(t, http.StatusForbidden, res.Code)
			} else {
				assert.Equal(t, http.StatusOK, res.Code)
				expectedRes := httptest.NewRecorder()
				encoder := json.NewEncoder(expectedRes)
				err:=encoder.Encode(os.Environ())
				assert.NoError(t,err,"couldn't encode to JSON format")
				assert.Equal(t, expectedRes.Body.String(), res.Body.String())
			}
		})
	}
}
