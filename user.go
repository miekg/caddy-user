package user

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(User{})
	httpcaddyfile.RegisterHandlerDirective("visitor_ip", parseCaddyfile)
}

// User holds the user id or username to we should use for serve requests.
type User struct {
	User   string `json:"user,omitempty"`
	uid    string
	gid    string
	logger *zap.Logger
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (u User) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return next.ServeHTTP(w, r)
}

// CaddyModule returns the Caddy module information.
func (User) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.user",
		New: func() caddy.Module { return new(User) },
	}
}

func (u *User) Provision(ctx caddy.Context) error {
	u.logger = ctx.Logger()
	return nil
}

func (u *User) Validate() error { return nil }

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (u *User) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume directive name

	if !d.NextArg() {
		return d.ArgErr()
	}

	u.User = d.Val() // lookup, uid and gid? TODO(miek)
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var u User
	err := u.UnmarshalCaddyfile(h.Dispenser)
	return u, err
}

var (
	_ caddy.Provisioner           = (*User)(nil)
	_ caddy.Validator             = (*User)(nil)
	_ caddyfile.Unmarshaler       = (*User)(nil)
	_ caddyhttp.MiddlewareHandler = (*User)(nil)
)
