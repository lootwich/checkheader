// Package checkheader provides a Traefik middleware plugin for checking group headers.
//
// This plugin checks the presence and validity of specific group headers in incoming
// HTTP requests. It is useful for authorization purposes, ensuring that requests
// contain the necessary group information in their headers.
//
// The plugin configuration requires the name of the header that contains the group
// information and a list of required groups. The header value is expected to be a
// comma-separated list of groups. The plugin will return an error response if the
// header is missing or if any of the required groups are not present in the header
// value.
//
// Example configuration:
//
//	{
//	  "groupHeaderName": "X-Group-Header",
//	  "neededGroups": ["admin", "user"]
//	}
//
// In this example, the plugin will look for the "X-Group-Header" header in incoming
// requests and check that it contains both the "admin" and "user" groups. If either
// group is missing, the request will be rejected with an "Unauthorized" error.
package checkheader

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	GroupHeaderName string   `json:"groupHeaderName,omitempty"`
	NeededGroups    []string `json:"neededGroups,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		GroupHeaderName: "",
		NeededGroups:    []string{},
	}
}

// CheckHeader a CheckHeader plugin.
type CheckHeader struct {
	next            http.Handler
	groupHeaderName string
	neededGroups    []string
	name            string
	template        *template.Template
}

// New created a new CheckHeader plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.GroupHeaderName == "" {
		return nil, fmt.Errorf("GroupHeaderName cannot be empty")
	}
	if len(config.NeededGroups) == 0 {
		return nil, fmt.Errorf("NeededGroups cannot be empty")
	}

	return &CheckHeader{
		next:            next,
		groupHeaderName: config.GroupHeaderName,
		neededGroups:    config.NeededGroups,
		name:            name,
		template:        template.New("demo").Delims("[[", "]]"),
	}, nil
}

func (a *CheckHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	headerValue := req.Header.Get(a.groupHeaderName)
	if headerValue == "" {
		http.Error(rw, "Missing required group header", http.StatusUnauthorized)
		return
	}

	// Split the header by comma and trim spaces
	headerGroups := make(map[string]struct{})
	for _, group := range strings.Split(headerValue, ",") {
		trimmed := strings.TrimSpace(group)
		if trimmed != "" {
			headerGroups[trimmed] = struct{}{}
		}
	}

	// Check if every needed group is present
	for _, needed := range a.neededGroups {
		if _, ok := headerGroups[needed]; !ok {
			http.Error(rw, "Unauthorized: missing required group", http.StatusUnauthorized)
			return
		}
	}

	a.next.ServeHTTP(rw, req)
}
