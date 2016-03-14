package main

import (
	"fmt"
	"io"
	"net/http"
)

func formfile(w http.ResponseWriter, r *http.Request) {
	key := "text file"
	file, hdr, err := r.FormFile(key)
	fmt.Println(file, hdr, err)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<form method="POST" enctype="multipart/form-data">
      <input type="file" name="text file">
      <input type="submit">
    </form>`)
}
func main() {
	http.HandleFunc("/", formfile)
	http.ListenAndServe(":8080", nil)

}
