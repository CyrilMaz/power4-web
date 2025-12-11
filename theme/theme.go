package theme

import (
	"net/http"
	"time"
)

const cookieName = "theme"
const defaultTheme = "light"

// GetTheme returns the theme stored in cookie ("light" or "dark").
// If no cookie or invalid value, returns the default theme.
func GetTheme(r *http.Request) string {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return defaultTheme
	}
	if c.Value != "light" && c.Value != "dark" {
		return defaultTheme
	}
	return c.Value
}

// SetTheme writes the theme cookie. Accepts only "light" or "dark".
func SetTheme(w http.ResponseWriter, theme string) {
	if theme != "light" && theme != "dark" {
		theme = defaultTheme
	}
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    theme,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
}

// ToggleHandler toggles the theme cookie or sets it when ?theme=light|dark is provided.
// Redirects back to the referer or to `/`.
func ToggleHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("theme")
	if q != "" {
		SetTheme(w, q)
	} else {
		cur := GetTheme(r)
		if cur == "dark" {
			SetTheme(w, "light")
		} else {
			SetTheme(w, "dark")
		}
	}

	ref := r.Referer()
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}
