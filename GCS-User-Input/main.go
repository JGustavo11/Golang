
package main

import (
	"golang.org/x/net/context"
  "google.golang.org/cloud/storage"
	"html/template"
  "log"
	"net/http"
	storageLog "google.golang.org/appengine/log"
	"io"
	"google.golang.org/appengine"

	
)

const BUCKET_NAME = "bucky"

func init() {
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css"))))
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/show", showHandler)
}

func userHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		file, header, err := req.FormFile("image")
		logError(err)
		userName := req.FormValue("userName")
		saveFile(req, userName, header.Filename, file)
		http.Redirect(res, req, "/show?userName="+userName, http.StatusFound)
		return
	}

	tpl := template.Must(template.ParseFiles("user.html"))
	err := tpl.Execute(res, nil)
	logError(err)
}

func saveFile(req *http.Request, userName string, fileName string, file io.Reader) {
	fileName = userName + "/" + fileName

	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	logStorageError(ctx, "Could not create a new client", err)
	defer client.Close()

	writer := client.Bucket(BUCKET_NAME).Object(fileName).NewWriter(ctx)
	writer.ACL = []storage.ACLRule{{
		storage.AllUsers,
		storage.RoleReader}}

	io.Copy(writer, file)
	writer.Close()
}

func showHandler(res http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)
	client, err := storage.NewClient(ctx)
	logStorageError(ctx, "Could not create a new client", err)
	defer client.Close()

	tpl := template.Must(template.ParseFiles("index.html"))
	err = tpl.Execute(res, PictureNames(ctx, client, inputUser(req)))
	logError(err)
}

func inputUser(req *http.Request) string {
	return req.FormValue("userName")
}

func PictureNames(ctx context.Context, client *storage.Client, userName string) []string {

	query := &storage.Query{
		Delimiter: "/",
		Prefix:    userName + "/",
	}
	objs, err := client.Bucket(BUCKET_NAME).List(ctx, query)
	logError(err)

	var names []string
	for _, result := range objs.Results {
		names = append(names, result.Name)
	}
	return names
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func logStorageError(ctx context.Context, errMessage string, err error) {
	if err != nil {
		storageLog.Errorf(ctx, errMessage, err)
	}
}
