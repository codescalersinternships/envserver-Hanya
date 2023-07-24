package envserver

import (
	"fmt"
	"io"
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
				assert.Nil(t, err)
				assert.NotNil(t, server)
				assert.Equal(t, tC.port, server.port)
			}
		})
	}
}

func TestServer_Run(t *testing.T) {
	t.Run("assert server runs correctly", func(t *testing.T) {
		port := 3000
		server, err := NewServer(port)
		assert.Nil(t, err)
		go func() {
			if err := server.Run(); err != nil {
				t.Errorf("Server.Run returned an error: %v", err)
			}
		}()
		resp, err := http.Get("http://localhost:" + fmt.Sprintf("%d", server.port) + "/env")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, err := readResponseBody(resp)
		assert.Nil(t, err)
		assert.Contains(t, body, "PATH=")
	})

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
			req := httptest.NewRequest("GET", "/env/"+tC.key, nil)
			writer := httptest.NewRecorder()
			envKeyHandler(writer, req)
			resp := writer.Result()
			body, err := readResponseBody(resp)
			body = strings.TrimSpace(body)
			if tC.returnsError {
				assert.Nil(t, err)
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Equal(t, value, body)
			}
		})
	}
}
func TestEnvHandler(t *testing.T) {
	t.Run("env variables are correct", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/env", nil)
		writer := httptest.NewRecorder()
		envHandler(writer, req)
		resp := writer.Result()
		body, err := readResponseBody(resp)
		body = strings.TrimSpace(body)
		value := os.Environ()
		strValue := strings.Join(value, "\n")
		assert.Nil(t, err)
		assert.NotNil(t, body)
		assert.Equal(t, strValue, body)
	})
}

func readResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	return string(bodyBytes), nil
}
