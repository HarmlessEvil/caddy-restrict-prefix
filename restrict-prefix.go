package restrictprefix

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func init() {
	caddy.RegisterModule(new(RestrictPrefix))
}

var (
	_ caddy.Provisioner           = (*RestrictPrefix)(nil)
	_ caddy.Validator             = (*RestrictPrefix)(nil)
	_ caddyhttp.MiddlewareHandler = (*RestrictPrefix)(nil)
)

type RestrictPrefix struct {
	Prefix string `json:"prefix,omitempty"`
	logger *zap.Logger
}

func (*RestrictPrefix) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.handlers.restrict_prefix",
		New: func() caddy.Module {
			return new(RestrictPrefix)
		},
	}
}

func (p *RestrictPrefix) Provision(ctx caddy.Context) error {
	p.logger = ctx.Logger()
	return nil
}

func (p *RestrictPrefix) Validate() error {
	if p.Prefix == "" {
		p.Prefix = "."
	}

	return nil
}

func (p *RestrictPrefix) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	for _, part := range strings.Split(r.URL.Path, "/") {
		if strings.HasPrefix(part, p.Prefix) {
			http.Error(w, "Not Found", http.StatusNotFound)

			if p.logger != nil {
				p.logger.Debug(fmt.Sprintf("restricted prefix: %q in %q", part, r.URL.Path))
			}

			return nil
		}
	}

	return next.ServeHTTP(w, r)
}
