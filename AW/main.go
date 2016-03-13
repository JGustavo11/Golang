package main
import (
	"net/http"
	"fmt"
)

func LeosOscar(res http.ResponseWriter,req *http.Request){
	fmt.Fprint(res,req.URL.Path)
}

func main(){
	http.HandleFunc("/",LeosOscar)
	http.ListenAndServe(":8080",nil)
}
