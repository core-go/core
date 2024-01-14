package cookies

import (
	"net/http"
	"time"
)

type Cookies struct {
	CookieName string
	Domain     string
	Expires    time.Duration
	SameSite   http.SameSite
}

func NewCookies(cookieName string, domain string, expires time.Duration, SameSite http.SameSite) *Cookies {
	return &Cookies{CookieName: cookieName, Domain: domain, Expires: expires, SameSite: SameSite}
}

func (c Cookies) RefreshValue(w http.ResponseWriter, sessionId string) error {
	http.SetCookie(w, &http.Cookie{
		Name:     c.CookieName,
		Domain:   c.Domain,
		Value:    sessionId,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   0,
		Expires:  time.Now().Add(c.Expires),
		SameSite: c.SameSite,
		Secure:   true,
	})
	return nil
}
