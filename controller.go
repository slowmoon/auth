package auth

import (
	"net/http"
	"path"
	"strings"

	"github.com/qor/auth/claims"
)

// NewServeMux generate http.Handler for auth
func (auth *Auth) NewServeMux() http.Handler {
	return &serveMux{Auth: auth}
}

type serveMux struct {
	*Auth
}

// ServeHTTP dispatches the handler registered in the matched route
func (serveMux *serveMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		claims  *claims.Claims
		reqPath = strings.TrimPrefix(req.URL.Path, serveMux.Prefix)
		paths   = strings.Split(reqPath, "/")
		context = &Context{Auth: serveMux.Auth, Claims: claims, Request: req, Writer: w}
	)

	if len(paths) >= 2 {
		// eg: /phone/login

		if provider := serveMux.Auth.GetProvider(paths[0]); provider != nil {
			context.Provider = provider

			// serve mux
			switch paths[1] {
			case "login":
				provider.Login(context)
			case "logout":
				provider.Logout(context)
			case "register":
				provider.Register(context)
			case "callback":
				provider.Callback(context)
			default:
				provider.ServeHTTP(context)
			}
			return
		}
	} else if len(paths) == 1 {
		// eg: /login, /logout

		switch paths[0] {
		case "login":
			// render login page
			serveMux.Auth.Render.Execute("auth/login", context, req, w)
			return
		case "logout":
			// destroy login context
			serveMux.Auth.Logout(w, req)
			return
		case "register":
			// render register page
			serveMux.Auth.Render.Execute("auth/register", context, req, w)
			return
		}
	}

	http.NotFound(w, req)
}

// AuthURL generate URL for auth
func (auth *Auth) AuthURL(pth string) string {
	return path.Join(auth.Prefix, pth)
}
