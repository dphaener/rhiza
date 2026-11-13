package web

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
)

// The web token turns the UI's actions on. TokenKey names the setting; the
// token's value is never shown: only whether it is set (TokenSet).
const (
	TokenKey = "rhiza.web.token"
	// TokenSet is what a surface that lists settings shows for a token that
	// is set, in place of its value.
	TokenSet = "set"
	// TokenCookie is the HttpOnly cookie a browser keeps the token in after
	// one visit with ?token=.
	TokenCookie = "rhiza_web_token"
	// CSRFField is the form field every action form carries; a POST that the
	// cookie authenticates is refused without the right value in it.
	CSRFField = "csrf"
	// MinTokenLen is the shortest token the setting accepts.
	MinTokenLen = 16
)

// ValidateToken checks a rhiza.web.token value.
func ValidateToken(v string) error {
	if len(v) < MinTokenLen || strings.IndexFunc(v, func(r rune) bool { return r <= ' ' || r == 0x7f }) >= 0 {
		return errTokenShape
	}
	return nil
}

type tokenError string

func (e tokenError) Error() string { return string(e) }

const errTokenShape = tokenError("must be at least 16 characters with no spaces or control characters")

type actionsKey struct{}

// ActionsOn reports whether the request carried the web token: a page shows
// its action controls only then.
func ActionsOn(r *http.Request) bool {
	on, _ := r.Context().Value(actionsKey{}).(bool)
	return on
}

// CSRF is the value a page puts in its action forms' CSRFField, or "" when
// the request did not carry the token (the page shows no actions then). It
// derives from the token, so it needs no state and is the one the gate accepts.
func CSRF(r *http.Request) string {
	v, _ := r.Context().Value(csrfKey{}).(string)
	return v
}

type csrfKey struct{}

func csrfFor(token string) string {
	m := hmac.New(sha256.New, []byte(token))
	m.Write([]byte("rhiza web csrf"))
	return hex.EncodeToString(m.Sum(nil))
}

// WithToken gates h's actions on the web token that token reads on each
// request (empty means none is set). Reads pass as they are. With no token
// set every request that is not a GET or HEAD is refused; with one, such a
// request needs it as an "Authorization: Bearer" header or in the HttpOnly
// cookie, or it is refused too. A POST the cookie authenticates (a browser
// sends it unasked) must also carry the CSRFField form field CSRF hands pages,
// or it is refused with 403; a bearer header is not ambient and needs none. A GET or HEAD with ?token=<token> sets the
// cookie and redirects to the same URL without it.
func WithToken(token func() string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		want := token()
		ok := want != "" && presented(r, want)
		bearer := want != "" && presentedBearer(r, want)
		read := r.Method == http.MethodGet || r.Method == http.MethodHead
		if read {
			if q := r.URL.Query(); want != "" && q.Has("token") {
				if equal(q.Get("token"), want) {
					http.SetCookie(w, &http.Cookie{Name: TokenCookie, Value: want, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode})
					q.Del("token")
					to := *r.URL
					to.RawQuery = q.Encode()
					w.Header().Set("Cache-Control", "no-store")
					http.Redirect(w, r, to.RequestURI(), http.StatusSeeOther)
					return
				}
			}
		} else if !ok {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			if want == "" {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("actions are off: no web token is set\n"))
				return
			}
			w.Header().Set("WWW-Authenticate", "Bearer")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("a valid web token is required\n"))
			return
		}
		if ok {
			csrf := csrfFor(want)
			if !read && !bearer && !equal(r.PostFormValue(CSRFField), csrf) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("a valid " + CSRFField + " field is required\n"))
				return
			}
			ctx := context.WithValue(r.Context(), actionsKey{}, true)
			r = r.WithContext(context.WithValue(ctx, csrfKey{}, csrf))
		}
		h.ServeHTTP(w, r)
	})
}

// presented reports whether r carries want, as a bearer header or the cookie.
func presented(r *http.Request, want string) bool {
	if presentedBearer(r, want) {
		return true
	}
	if c, err := r.Cookie(TokenCookie); err == nil && equal(c.Value, want) {
		return true
	}
	return false
}

func equal(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func presentedBearer(r *http.Request, want string) bool {
	a, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	return found && equal(strings.TrimSpace(a), want)
}
