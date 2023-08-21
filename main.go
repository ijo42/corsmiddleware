package corsmiddleware

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var DefaultAllowHeaders []string = []string{"Content-Type", "Content-Length",
	"Accept-Encoding", "Authorization", "Accept", "Origin", "Referer", "Cache-Control"}

type Config struct {
	AllowCredentials bool     `json:"allow_credentials,omitempty"`
	AllowOrigins     []string `json:"allow_origins,omitempty"`
	AllowMethods     []string `json:"allow_methods,omitempty"`
	AllowHeaders     []string `json:"allow_headers,omitempty"`
	ExposeHeaders    []string `json:"expose_headers,omitempty"`
	MaxAge           int64    `json:"max_age,omitempty"`
}

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

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	var err error

	var origins []*regexp.Regexp

	if Contains(config.AllowOrigins, "*") {
		all, _ := regexp.Compile(".*")
		origins = []*regexp.Regexp{all}
	} else {
		origins, err = CompileOrigins(config.AllowOrigins)
		if err != nil {
			return nil, err
		}
	}

	headers := MergeAndUniques(DefaultAllowHeaders, config.AllowHeaders)

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

	if req.Method != "OPTIONS" {
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
