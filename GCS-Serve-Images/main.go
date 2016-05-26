package main

import (
	"golang.org/x/net/context"
    storageLog "google.golang.org/appengine/log"
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/cloud/storage"
	"html/template"
	"log"
)

const BUCKET_NAME = "bucky"

func init() {
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.HandleFunc("/", handler)
}

func handler(res http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	StoreLogError(ctx, "unable to create", err)
	defer client.Close()

	tpl := template.Must(template.ParseFiles("index.html"))
	err = tpl.Execute(res, picNames(ctx, client))
	logError(err)
}

func StoreLogError(ctx context.Context, errMessage string, err error) {
	if err != nil {
		storageLog.Errorf(ctx, errMessage, err)
	}
}
func picNames(ctx context.Context, client *storage.Client) []string {

	query := &storage.Query{
		Delimiter: "/",
		Prefix:    "photos/",
	}
	objs, err := client.Bucket(BUCKET_NAME).List(ctx, query)
	logError(err)

	var nombre []string
	for _, result := range objs.Results {
		nombre = append(nombre, result.Name)
	}
	return nombre
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}
