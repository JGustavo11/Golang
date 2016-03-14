package main

import (
	"fmt"
	"net/http"
	"github.com/nu7hatch/gouuid"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		cook, err := req.Cookie("my-myyy")
		if err != nil {
			id, _ := uuid.New()
			cook = &http.Cookie{
				Name:  "my-session",
				Value: id.String(),

				HttpOnly: true,
			}
			http.SetCookie(res, cook)
		}
		fmt.Printf("%v", cook)
	})
	http.ListenAndServe(":8080", nil)
}
