package main

import (
	"log"
	"os"
	"text/template"
)
type person struct{
	Name string ///because string
}

type conditions struct {
	person
	torf bool

}

func main() {

	p1 := conditions{
		person: person{
			Name: "ScottyRace" ,
		},
		torf:false,
	}

	if p1.Name == "scotty" {
		p1.torf = true
	}

	tmp, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatalln(err)
	}

	err = tmp.Execute(os.Stdout, p1)
	if err != nil {
		log.Fatalln(err)
	}
}