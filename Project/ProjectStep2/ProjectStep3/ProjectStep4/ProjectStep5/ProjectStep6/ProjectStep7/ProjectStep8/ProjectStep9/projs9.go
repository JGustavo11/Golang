package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/nu7hatch/gouuid"
//"html/template"
	"io/ioutil"
	"net/http"
	"io"
	"html/template"
)

type User struct {
	Uuid, Name, Age, Hmac string
}


func main() {
	http.HandleFunc("/", serveLogin)
	http.ListenAndServe(":8080", nil)
}
var loginFile string

func init() {
	temp, _ := ioutil.ReadFile("user_templates/temp2.html")
	loginFile = string(temp)
}

func getHMAC(data string) string {
	hmac_code := hmac.New(sha256.New, []byte(data+"plebKey"))
	return string(hmac_code.Sum(nil))
}
func bakeUserCookie(cookie *http.Cookie, req *http.Request) string {
	jsonVal, _ := undoJSON(cookie)
	jsonVal.Name = req.FormValue("name")
	jsonVal.Age = req.FormValue("age")
	jsonVal.Hmac = req.FormValue("HMAC")
	return redoJSON(jsonVal)
}

func redoJSON(jsonVal User) string {
	encode, _ := json.Marshal(jsonVal)
	return base64.StdEncoding.EncodeToString(encode)
}

func undoJSON(cookie *http.Cookie) (User, bool) {
	decode, _ := base64.StdEncoding.DecodeString(cookie.Value)
	var jsonVal User
	json.Unmarshal(decode, &jsonVal)
	if hmac.Equal([]byte(jsonVal.Hmac), []byte(getHMAC(jsonVal.Uuid+jsonVal.Name+jsonVal.Age))) {
		return jsonVal, true
	}
	return jsonVal, false
}


func serveLogin(res http.ResponseWriter, req *http.Request) {

	cookie, err := req.Cookie("logged-in")

	if err == http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:  "logged-in",
			Value: "0",
			//Secure: true,
			HttpOnly: true,
		}
	}

	// check log in
	if req.Method == "POST" {
		password := req.FormValue("password")
		if password == "pleb" {
			cookie = &http.Cookie{
				Name:  "logged-in",
				Value: "1",
				//Secure: true,
				HttpOnly: true,
			}


		}
	}

	// if logout, then logout
	if req.URL.Path == "/logout" {
		cookie = &http.Cookie{
			Name:   "logged-in",
			Value:  "0",
			MaxAge: -1,
			//Secure: true,
			HttpOnly: true,
		}
		http.SetCookie(res, cookie)
		http.Redirect(res, req, "/", 303)
		return
	}

	http.SetCookie(res, cookie)
	var html string

	// not logged in
	if cookie.Value == "0" {
		html = `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title></title>
			</head>
			<body>
			<h1>LOG IN</h1>
			<form method="POST">
				<h3>Password</h3>
				<input type="text" name="password">
				<br>
				<input type="submit">
			</form>
			</body>
			</html>`
	}

	// logged in
	if cookie.Value == "1" {
		html = `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title></title>
			</head>
			<body>
			<h1><a href="/logout">LOG OUT</a></h1>
			</body>
			</html>`
		serveUserData(res,req)
		bakeUserCookie(cookie, req)
	}

	io.WriteString(res, html)
}

func serveUserData(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session-fino")
	if err != nil {
		cookie = userCookie()
		http.SetCookie(res, cookie)
	}
	if req.Method == "POST" {
		cookie.Value = bakeUserCookie(cookie, req)
	}
	obj, _ := undoJSON(cookie)

	t, _ := template.New("Name").Parse(loginFile)
	t.Execute(res, obj)
}


func userCookie() *http.Cookie {
	id, _ := uuid.NewV4()
	temp := User{Uuid: id.String(), Hmac: getHMAC(id.String())}
	b, _ := json.Marshal(temp)
	encode := base64.StdEncoding.EncodeToString(b)
	return &http.Cookie{
		Name:     "session-fino",
		Value:    encode,
		HttpOnly: true,
		//Secure : true,
	}
}
