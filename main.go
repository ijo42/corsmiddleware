// Package corsmiddleware plugin.
package corsmiddleware

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var defaultAllowHeaders = []string{
	"Content-Type", "Content-Length",
	"Accept-Encoding", "Authorization", "Accept", "Origin", "Referer", "Cache-Control",
}

// Config the plugin configuration.
type Config struct {
	AllowCredentials bool     `json:"allowCredentials,omitempty"`
	AllowOrigins     []string `json:"allowOrigins,omitempty"`
	AllowMethods     []string `json:"allowMethods,omitempty"`
	AllowHeaders     []string `json:"allowHeaders,omitempty"`
	ExposeHeaders    []string `json:"exposeHeaders,omitempty"`
	MaxAge           int64    `json:"maxAge,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		AllowCredentials: false,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{},
		ExposeHeaders:    []string{},
		MaxAge:           86400,
	}
}

// CORS structure for plugin execution.
type CORS struct {
	next        http.Handler
	name        string
	credentials bool

	origins []*regexp.Regexp
	methods []string
	headers []string
	exposed []string

	age int64
}

// New created a new Demo plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	var err error

	var origins []*regexp.Regexp

	if Contains(config.AllowOrigins, "*") {
		all := regexp.MustCompile(".*")
		origins = []*regexp.Regexp{all}
	} else {
		origins, err = CompileOrigins(config.AllowOrigins)
		if err != nil {
			return nil, err
		}
	}

	headers := MergeAndUniques(defaultAllowHeaders, config.AllowHeaders)

	return &CORS{
		next: next,
		name: name,

		origins:     origins,
		headers:     headers,
		methods:     config.AllowMethods,
		exposed:     config.ExposeHeaders,
		age:         config.MaxAge,
		credentials: config.AllowCredentials,
	}, nil
}

func (c *CORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodOptions {
		c.next.ServeHTTP(res, req)
		return
	}

	origin := req.Header.Get("Origin")
	if origin == "" {
		http.Error(res, "No `Origin` header received", http.StatusBadRequest)
		return
	}

	if !AllowOrigin(c.origins, origin) {
		// Response is sent without headers as cors is not allowed.
		res.WriteHeader(http.StatusNoContent)
	}

	res.Header().Set("Access-Control-Allow-Origin", origin)
	res.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(c.credentials))
	res.Header().Set("Access-Control-Allow-Headers", strings.Join(c.headers, ", "))
	res.Header().Set("Access-Control-Allow-Methods", strings.Join(c.methods, ", "))
	res.Header().Set("Access-Control-Max-Age", strconv.FormatInt(c.age, 10))
	res.WriteHeader(http.StatusNoContent)
}
