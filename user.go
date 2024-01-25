package user

import (
	"log"
	"net/http"
	"os/user"
	"runtime"
	"strconv"
	"syscall"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(&User{})
	httpcaddyfile.RegisterHandlerDirective("user", parseCaddyfile)
}

// User holds the user id or username to we should use for serve requests.
type User struct {
	User string `json:"user,omitempty"`
	Uid  uintptr
	l    *zap.Logger
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (u *User) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// most of this stolen from: https://stackoverflow.com/questions/56403237/is-it-possible-to-run-a-goroutine-or-go-method-under-a-different-user
	// TL;DR: we lock the goroutine to a thread and then call setuid
	runtime.LockOSThread()
	u.l.Sugar().Infof("Setting to %s (%d)", u.User, u.Uid)
	if _, _, errno := syscall.Syscall(syscall.SYS_SETUID, u.Uid, 0, 0); errno != 0 {
		err := errno
		u.l.Sugar().Warnf("Unable to set user to: %s: %d: %s", u.User, errno, err)
	}
	u.l.Sugar().Infof("uid: %d", syscall.Getuid())
	u.l.Sugar().Infof("euid: %d", syscall.Geteuid())
	err := next.ServeHTTP(w, r)
	runtime.UnlockOSThread()
	return err
}

// CaddyModule returns the Caddy module information.
func (u *User) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.user",
		New: func() caddy.Module { return new(User) },
	}
}

func (u *User) Provision(ctx caddy.Context) error {
	u.l = ctx.Logger()
	return nil
}

func (u *User) Validate() error { return nil }

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (u *User) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume directive name

	if !d.NextArg() {
		return d.ArgErr()
	}

	u.User = d.Val()
	u1, err := user.Lookup(u.User)
	if err != nil {
		return err
	}
	log.Printf("User: " + u.User)
	uid, err := strconv.ParseUint(u1.Uid, 10, 64)
	if err != nil {
		return err
	}
	log.Printf("Uid: %d", uid)
	u.Uid = uintptr(uid)

	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	u := &User{}
	err := u.UnmarshalCaddyfile(h.Dispenser)
	return u, err
}

var (
	_ caddy.Provisioner           = (*User)(nil)
	_ caddy.Validator             = (*User)(nil)
	_ caddyfile.Unmarshaler       = (*User)(nil)
	_ caddyhttp.MiddlewareHandler = (*User)(nil)
)
