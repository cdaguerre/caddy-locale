package locale

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"golang.org/x/text/language"
)

func Test(t *testing.T) {
	for i, tc := range []struct {
		test string
		available []string
		config Handler
		header string
		expect string
	}{
		{
			test: "Minimal",
			available: []string{"en"},
			config: Handler{
				Methods: []string{"header", "cookie"},
			},
			header: "fr",
			expect: "en",
		},
		{
			test: "Not default locale",
			available: []string{"fr", "en"},
			config: Handler{
				Methods: []string{"header"},
			},
			header: "en",
			expect: "en",
		},		
		{
			test: "First locale is default fallback",
			available: []string{"fr", "en"},
			config: Handler{
				Methods: []string{"header"},
			},
			header: "es",
			expect: "fr",
		},				
		{
			test: "With regional locale",
			available: []string{"en", "en-US"},
			config: Handler{
				Methods: []string{"header"},
			},
			header: "en-US",
			expect: "en-US",
		},	
		{
			test: "Header with weights",
			available: []string{"en", "en-US", "fr"},
			config: Handler{
				Methods: []string{"header"},
			},
			header: "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7",
			expect: "fr",
		},						
		{
			test: "Falls back to locale without region",
			available: []string{"fr", "en"},
			config: Handler{
				Methods: []string{"header"},
			},
			header: "en-US",
			expect: "en",
		},	
		{
			test: "Case sensitivity",
			available: []string{"en", "en-us", "fr"},
			config: Handler{
				Methods: []string{"header"},
			},
			header: "en-us",
			expect: "en-US",
		},				
	} {
		h := tc.config

		for _, locale := range tc.available {
			tag := language.Make(locale)

			h.AvailableLocales = append(h.AvailableLocales, tag)
		}		

		req, _ := http.NewRequest(http.MethodGet, "http://localhost:9080/", nil)
		req.Header.Set("accept-language", tc.header)
		w := httptest.NewRecorder()
	
		var handler caddyhttp.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) error {
			actual := request.Header.Get("Detected-Locale")
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