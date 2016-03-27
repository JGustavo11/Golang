package main

import (

	"fmt"
	"github.com/nu7hatch/gouuid"
	"html/template"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Name string
	Age  string
}

func main() {
	tpl, err := template.ParseFiles("index4.html")
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		name := req.FormValue("name")
		age := req.FormValue("age")
		data := foo(name, age)
		cookie, err := req.Cookie("session-fino")
		if err != nil {
			id, _ := uuid.NewV4()
			cookie = &http.Cookie{
				Name:  "session-fino",
				Value: id.String() + "|" + name + age + "|" + data,
				// Secure: true,
				HttpOnly: true,
			}
			http.SetCookie(res, cookie)
		}
		err = tpl.Execute(res, nil)
		if err != nil {
			log.Fatalln(err)
		}
	})
	http.ListenAndServe(":8080", nil)
}

func foo(name string, age string) string {
	user := User{
		Name: name,
		Age:  age,
	}

	bs, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error occured !!!: ", err)
	}
	str := base64.URLEncoding.EncodeToString(bs)
	return str
}
