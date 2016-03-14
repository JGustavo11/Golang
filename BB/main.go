package main

import (
	"io"
	"net/http"
	"strconv"
)
func cookiemonster(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	cookie, err := r.Cookie("Cvalue")
	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "Cvalue",
			Value: "0",
		}
	}

	count, _ := strconv.Atoi(cookie.Value)
	count++
	cookie.Value = strconv.Itoa(count)

	http.SetCookie(w, cookie)
	

	io.WriteString(w,"value:" +cookie.Value)
}
func main() {
	http.HandleFunc("/", cookiemonster)
	http.ListenAndServe(":9000", nil)

}
