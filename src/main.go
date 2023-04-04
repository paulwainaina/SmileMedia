package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)
var (
	tpl   *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("./templates/*.html"))
}

type Page struct {
	Body  []byte
	Title string
	Data  interface{}
}

func LoadPage(file string) (*Page, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var body []byte
	_, err = f.Read(body)
	if err != nil {
		return nil, err
	}
	return &Page{Body: body}, nil
}

func RenderTemplate(w http.ResponseWriter, file string, page *Page) {
	err := tpl.ExecuteTemplate(w, file, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func IndexHandler(w http.ResponseWriter,r *http.Request ){
	file := "index.html"
	filePath := "templates/" + file
	pageName := "Smile media Home"
	page, err := LoadPage(filePath)
	if err != nil {
		page = &Page{Title: pageName}
	}
	page.Data = nil
	RenderTemplate(w, file, page)
}

func AboutHandler(w http.ResponseWriter,r *http.Request ){
	file := "about.html"
	filePath := "templates/" + file
	pageName := "Smile media Contact"
	page, err := LoadPage(filePath)
	if err != nil {
		page = &Page{Title: pageName}
	}
	page.Data = nil
	RenderTemplate(w, file, page)
}

func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("internal server error loading environment")
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.Handle("/index", http.HandlerFunc(IndexHandler))
	http.Handle("/about", http.HandlerFunc(AboutHandler))
	http.Handle("/", http.RedirectHandler("/index", http.StatusSeeOther))
	err = http.ListenAndServe(string(os.Getenv("SERVER")+":"+os.Getenv("PORT")), nil)
	if err == http.ErrServerClosed {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Println("server error occured : " + err.Error())
		os.Exit(1)
	}
}