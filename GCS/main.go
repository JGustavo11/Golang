package main

import (
	"net/http"
	"google.golang.org/appengine/log"
	"io"
	"google.golang.org/appengine"
)



func init() {
	http.HandleFunc("/", baseS)
	http.HandleFunc("/get", get)
	http.Handle("/favicon.ico", http.NotFoundHandler())
}


func baseS(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	html := `
		<h1>Upload File</h1>
	    <form method="POST" enctype="multipart/form-data">
			<input type="file" name="data">
			<input type="submit">
	    </form>
	`

	if req.Method == "POST" {
		mpf, hdr, err := req.FormFile("data")
		
		if err != nil {
			log.Errorf(ctx, "*** ctx ***", err)
			http.Error(res, "*** Unable to upload file ***", http.StatusInternalServerError)
			return
		}
		defer mpf.Close()

	
		filename, err := storeFile(req, mpf, hdr)
		if err != nil {
			log.Errorf(ctx, "*** ***")
			http.Error(res, "***  unable to accept file ***\n" + err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		filenames, err := genCookie(res, req, filename)
		if err != nil {
			log.Errorf(ctx, "*** CTX ERROR: ", err, "***")
			http.Error(res, "*** ERROR:  unable to accept file\n" + err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		html += `<h1>Files</h1>`
		for f, _ := range filenames{
			html += `<a href="/get?object=` + f + `">` + f + `</a><br>`
		}
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(res, html)
}


func get(res http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	object := req.FormValue("object")
	rdr, err := get(ctx, object)
	if err != nil {
		log.Errorf(ctx, "*** CTX ERROR: ***", err)
		http.Error(res, "*** HTTP ERROR:  unable to get file " + object + " ***\n" + err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	defer rdr.Close()
	io.Copy(res, rdr)
}
