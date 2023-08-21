package corsmiddleware_test

import (
	"context"
	"fmt"
	"gotest.tools/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SergioFloresG/corsmiddleware"
)

func TestStaticDomainSuccess(t *testing.T) {
	var testCases = []struct {
		domain string
		origin string
	}{
		// Unsecure
		{domain: "http://localhost", origin: "http://localhost"},
		{domain: "http://example.com", origin: "http://example.com"},
		{domain: "http://example.com", origin: "http://example.com"},
		{domain: "http://internal.example.com", origin: "http://internal.example.com"},
		// Secure
		{domain: "https://example.com", origin: "https://example.com"},
		{domain: "https://example.com", origin: "https://example.com"},
		{domain: "https://internal.example.com", origin: "https://internal.example.com"},
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	var req *http.Request
	var conf *corsmiddleware.Config
	var err error

	setupTest := func(tb testing.TB) {
		conf = corsmiddleware.CreateConfig()
		req, err = http.NewRequestWithContext(ctx, http.MethodOptions, "https://localhost", nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("[%s]", tc.domain), func(st *testing.T) {
			setupTest(t)
			recorder := httptest.NewRecorder()

			conf.AllowOrigins = []string{tc.domain}
			req.Header.Set("Origin", tc.origin)

			handler, err := corsmiddleware.New(ctx, next, conf, "statics")
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(recorder, req)

			assert.Equal(t, recorder.Code, http.StatusNoContent)
			assert.Equal(t, recorder.Header().Get("Access-Control-Allow-Origin"), tc.domain)
		})
	}

}

func TestWildcardDomainSuccess(t *testing.T) {
	var allowOrigins = []string{"https://*.foo.com", "https://*.example.com", "https://localhost"}
	var testCases = []string{
		"https://localhost",
		"https://bar.foo.com",
		"https://loop.foo.com",
		"https://netpay.foo.com",
		"https://sub.example.com",
		"https://foo.example.com",
		"https://bar.example.com",
	}

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	var req *http.Request
	var err error

	conf := corsmiddleware.CreateConfig()
	conf.AllowOrigins = allowOrigins

	setupTest := func(tb testing.TB) {
		req, err = http.NewRequestWithContext(ctx, http.MethodOptions, "http://localhost", nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	testHelper := func(t *testing.T, originDomain string) {
		t.Helper()
		recorder := httptest.NewRecorder()
		req.Header.Set("Origin", originDomain)

		handler, err := corsmiddleware.New(ctx, next, conf, "wildcard")
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(recorder, req)

		assert.Equal(t, recorder.Code, http.StatusNoContent)
		assert.Equal(t, recorder.Header().Get("Access-Control-Allow-Origin"), originDomain)
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("[%s]", tc), func(t *testing.T) {
			setupTest(t)
			testHelper(t, tc)
		})
	}

	conf.AllowOrigins = []string{"*"}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("AnyDomain[%s]", tc), func(t *testing.T) {
			setupTest(t)
			testHelper(t, tc)
		})
	}

}
