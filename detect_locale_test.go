package locale

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddytest"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"golang.org/x/text/language"
)

func TestConfig(t *testing.T) {
	tester := caddytest.NewTester(t)
	tester.InitServer(`
	{
		admin localhost:2999
		http_port     9080
		https_port    9443
	}
	localhost:9080 {
		route / {
			locale en fr de
			respond {vars.detected-locale}
		}
	}`, "caddyfile")

	req1, _ := http.NewRequest(http.MethodGet, "http://localhost:9080/", nil)
	req1.Header.Set("accept-language", "de")
	tester.AssertResponse(req1, 200, "de")

	req2, _ := http.NewRequest(http.MethodGet, "http://localhost:9080/", nil)
	req2.Header.Set("accept-language", "es")
	tester.AssertResponse(req2, 200, "en")	
}

func TestCases(t *testing.T) {
	for i, tc := range []struct {
		test string
		available []string
		header string
		expect string
	}{
		{
			test: "Minimal",
			available: []string{"en"},
			header: "fr",
			expect: "en",
		},
		{
			test: "Not default locale",
			available: []string{"fr", "en"},
			header: "en",
			expect: "en",
		},		
		{
			test: "First locale is default fallback",
			available: []string{"fr", "en"},
			header: "es",
			expect: "fr",
		},				
		{
			test: "With regional locale",
			available: []string{"en", "en-US"},
			header: "en-US",
			expect: "en-US",
		},	
		{
			test: "Header with weights",
			available: []string{"en", "en-US", "fr"},
			header: "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
			expect: "fr",
		},						
		{
			test: "Falls back to locale without region",
			available: []string{"fr", "en"},
			header: "en-US",
			expect: "en",
		},	
		{
			test: "Case sensitivity",
			available: []string{"en", "en-us", "fr"},
			header: "en-us",
			expect: "en-US",
		},				
	} {
		h := DetectLocale{}
		h.Methods = []string{"header"}
		h.HeaderName = "Accept-Language"

		for _, locale := range tc.available {
			tag := language.Make(locale)

			h.AvailableLocales = append(h.AvailableLocales, tag)
		}		

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:9080/", nil)
		req.Header.Set("accept-language", tc.header)
		w := httptest.NewRecorder()
	
		var handler caddyhttp.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) error {
			actual := request.Header.Get(h.HeaderName)
			if actual == "" {
				t.Errorf("Test %d (%s): Empty header", i, tc.test)
			}
			if actual != tc.expect {
				t.Errorf("Test %d (%s): Expected=%s Actual=%s", i, tc.test, tc.expect, actual)
			}			
	
			return nil
		}
	
		_, cancel := caddy.NewContext(caddy.Context{Context: context.Background()})
		defer cancel()
	
	
		if err := h.ServeHTTP(w, req, handler); err != nil {
			t.Errorf("Test %d (%s): ServeHTTP error: %v", i, tc.test, err)
			continue
		}
	}
}