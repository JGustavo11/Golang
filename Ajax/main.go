package main

import (
	"google.golang.org/appengine/log"
    "html/template"
	"google.golang.org/appengine/datastore"
	"net/http"
    "google.golang.org/appengine"
	"encoding/json"
	
)

type Word struct {
	Name string
}

var tpl *template.Template

func init() {
	http.HandleFunc("/", theIndex)
	http.HandleFunc("/api/check", checkWord)

	
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))


	tpl = template.Must(template.ParseGlob("*.html"))
}

func theIndex(res http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {

		var w Word
		w.Name = req.FormValue("new-word")

		ctx := appengine.NewContext(req)
		log.Infof(ctx, "submitted", w.Name)

		key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
		_, err := datastore.Put(ctx, key, &w)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tpl.ExecuteTemplate(res, "theIndex.html", nil)
}

func checkWord(res http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)


	var w Word
	json.NewDecoder(req.Body).Decode(&w)
	log.Infof(ctx, "ENTERED checkWord - w.Name: %v", w.Name)


	key := datastore.NewKey(ctx, "Dictionary", w.Name, 0, nil)
	err := datastore.Get(ctx, key, &w)
	if err != nil {
		json.NewEncoder(res).Encode("false")
		return
	}
	json.NewEncoder(res).Encode("true")
}
