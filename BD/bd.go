package main

import (
	"fmt"
	"io"
	"net/http"
	"crypto/hmac"
	"crypto/sha256"

)

func getCode(data string) string {
	h := hmac.New(sha256.New, []byte("ourkey"))
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func foo(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	cookie, err := r.Cookie("session-id")
	if err != nil {
		cookie = &http.Cookie{
			Name:     "session-id",
			Value:    "cookie-value",
			HttpOnly: true,
		}
	}

	if r.FormValue("email") != "" {
		email := r.FormValue("email")
		cookie.Value = email + `|` + getCode(email)
	}

	http.SetCookie(w, cookie)
	io.WriteString(w, `<!DOCTYPE html>
	<html>
	  <body>
	    <form method="POST">
	    `+cookie.Value+`
	      <input type="email" name="email">
	      <input type="submit">
	    </form>
	  </body>
	</html>`)

}
func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":8080", nil)
}
