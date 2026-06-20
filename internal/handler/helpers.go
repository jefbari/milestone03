package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// parseID extracts a named path param from URL path.
// Works with stdlib net/http path values (Go 1.22+).
func parseID(r *http.Request, name string) (int64, error) {
	val := r.PathValue(name)
	if val == "" {
		// fallback: try to grab last segment if PathValue not set
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) == 0 {
			return 0, errors.New("missing id")
		}
		val = parts[len(parts)-1]
	}
	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
