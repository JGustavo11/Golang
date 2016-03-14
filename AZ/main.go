package main

import (
	"net/http"
	"io"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		key := "name"
		value := req.FormValue(key)
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		if value != "" {
			value = "Name: " + value;
			io.WriteString(res, value);
		} else {
			io.WriteString(res,
				`<form method = "GET">
				<input type = "text" name = "Como Te Llamas"/>
				<input type = "SUBMIT"/>
			</form>`)
		}
	})
	http.ListenAndServe(":9000", nil)
}
