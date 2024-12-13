package locale

import (
	"net/http"
	"slices"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"

	"golang.org/x/text/language"
)

const moduleName = "locale"

func init() {
	caddy.RegisterModule(Handler{})
	httpcaddyfile.RegisterHandlerDirective(moduleName, parseCaddyfileHandlerDirective)
	httpcaddyfile.RegisterDirectiveOrder(moduleName, httpcaddyfile.Before, "rewrite")
}

// Handler is a httpserver to detect the user's locale.
type Handler struct {
	AvailableLocales []language.Tag `json:"locales"`
	Methods          []string `json:"methods"`
	CookieName       string
	HeaderName       string
}

// CaddyModule returns the Caddy module information.
func (Handler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.detect_locale",
		New: func() caddy.Module { return new(Handler) },
	}
}

// ServeHTTP implements caddyhttp.HandlerHandler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	var tags = []language.Tag{}

	var matcher = language.NewMatcher(h.AvailableLocales)

	if slices.Contains(h.Methods, "cookie") {
		lang, _ := r.Cookie("lang")
		tags = append(tags, language.Make(lang.String()))
	}
	if slices.Contains(h.Methods, "header") {
		t, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		tags = append(tags, t...)
	}	
	
	tag, _, _ := matcher.Match(tags...)
	base, _, region := tag.Raw()

	locale := base.String()
	regionSuffix := region.String()

	if regionSuffix != "ZZ" {
		locale = locale + "-" + regionSuffix
	}

	r.Header.Set(h.HeaderName, locale)

	return next.ServeHTTP(w, r)
}


func parseCaddyfileHandlerDirective(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var hl Handler
	return &hl, hl.UnmarshalCaddyfile(h.Dispenser)
}

func (h *Handler) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		h.CookieName = "lang"
		h.HeaderName = "Detected-Locale"
		h.Methods = []string{"header"}

		localeArgs := d.RemainingArgs()
		for _, localeArg := range localeArgs {
			tag := language.Make(localeArg)

			h.AvailableLocales = append(h.AvailableLocales, tag)
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {

			case "available":
				localeArgs := d.RemainingArgs()
				for _, localeArg := range localeArgs {
					tag := language.Make(localeArg)

					h.AvailableLocales = append(h.AvailableLocales, tag)
				}
			case "methods":
				detectArgs := d.RemainingArgs()
				h.Methods = append([]string{}, detectArgs...)
			case "cookie":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if value := strings.TrimSpace(d.Val()); value != "" {
					h.CookieName = value
				}
			case "header":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if value := strings.TrimSpace(d.Val()); value != "" {
					h.HeaderName = value
				}				
			default:
				return d.ArgErr()
			}
		}		
	}

	return nil
}
