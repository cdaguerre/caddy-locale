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

func init() {
	caddy.RegisterModule(DetectLocale{})
	httpcaddyfile.RegisterHandlerDirective("locale", parseCaddyfileHandlerDirective)
	httpcaddyfile.RegisterDirectiveOrder("locale", httpcaddyfile.Before, "rewrite")
}

// Detect and normalize user locale.
//
// Syntax:
//
//	locale en de fr
type DetectLocale struct {
	AvailableLocales []language.Tag `json:"locales"`
	Methods          []string `json:"methods"`
	CookieName       string `json:"cookie"`
	HeaderName       string `json:"header"`
}

// CaddyModule returns the Caddy module information.
func (DetectLocale) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.detect_locale",
		New: func() caddy.Module { return new(DetectLocale) },
	}
}

// Provision implements caddy.Provisioner.
func (dl *DetectLocale) Provision(ctx caddy.Context) error {
	return nil
}

// Validate implements caddy.Validator.
func (dl *DetectLocale) Validate() error {
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (dl *DetectLocale) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	var tags = []language.Tag{}

	var matcher = language.NewMatcher(dl.AvailableLocales)

	if slices.Contains(dl.Methods, "cookie") {
		lang, _ := r.Cookie("lang")
		tags = append(tags, language.Make(lang.String()))
	}
	if slices.Contains(dl.Methods, "header") {
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

	if dl.HeaderName != "" {
		r.Header.Set(dl.HeaderName, locale)
	}
		
	caddyhttp.SetVar(r.Context(), "detected-locale", locale)

	return next.ServeHTTP(w, r)
}

func (dl *DetectLocale) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	dl.CookieName = "lang"
	dl.HeaderName = "Detected-Locale"
	dl.Methods = []string{"header"}

	for d.Next() {
		localeArgs := d.RemainingArgs()
		for _, localeArg := range localeArgs {
			tag := language.Make(localeArg)

			dl.AvailableLocales = append(dl.AvailableLocales, tag)
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "available":
				localeArgs := d.RemainingArgs()
				for _, localeArg := range localeArgs {
					tag := language.Make(localeArg)
					dl.AvailableLocales = append(dl.AvailableLocales, tag)
				}
			case "methods":
				detectArgs := d.RemainingArgs()
				dl.Methods = append([]string{}, detectArgs...)
			case "cookie":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if value := strings.TrimSpace(d.Val()); value != "" {
					dl.CookieName = value
				}
			case "header":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if value := strings.TrimSpace(d.Val()); value != "" {
					dl.HeaderName = value
				}				
			default:
				return d.ArgErr()
			}
		}		
	}

	return nil
}

func parseCaddyfileHandlerDirective(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var dl DetectLocale
	return &dl, dl.UnmarshalCaddyfile(h.Dispenser)
}

// Interface guards
var (
	_ caddy.Provisioner           = (*DetectLocale)(nil)
	_ caddy.Validator             = (*DetectLocale)(nil)
	_ caddyhttp.MiddlewareHandler = (*DetectLocale)(nil)
	_ caddyfile.Unmarshaler       = (*DetectLocale)(nil)
)
