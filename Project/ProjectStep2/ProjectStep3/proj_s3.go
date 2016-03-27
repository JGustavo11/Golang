package main

import (
//"fmt"
	"github.com/nu7hatch/gouuid"//configured
	"html/template"
	"log"
	"net/http"
)

func main() {
	tpl, err := template.ParseFiles("index3.html")
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		name := req.FormValue("Name:")
		age := req.FormValue("Age:")
		cookie, err := req.Cookie("session-fino")
		if err != nil {
			id, _ := uuid.NewV4()
			cookie = &http.Cookie{
				Name:  "session-fino",
				Value: id.String() + "|" + name + age,
				HttpOnly: true,
			}
			http.SetCookie(res, cookie)
		}
		err = tpl.Execute(res, nil)
		if err != nil {
			log.Fatalln(err)
		}
	})
	http.ListenAndServe(":9000", nil)
}
