package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates *template.Template
var templatesDir = "templates"
var allowedKinds = []string{"image/png", "image/jpeg"}

func init() {
	templates = parseTemplates(templatesDir, ".gohtml")
}
func parseTemplates(path, extension string) *template.Template {
	templ := template.New("")
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, extension) {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	return templ
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	//Parse our multipart form, 10 << 20 specifies a maximum
	//upload of 10 MB files.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get a ref to the parsed multipart form
	m := r.MultipartForm
	fmt.Println(m)

	//get the *fileheaders
	files := m.File["myFile"]
	fmt.Println(files)
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
		fmt.Println(file)
		if _, err = file.Read(buff); err != nil {
			fmt.Println(err) // do something with that error
			return
		}
		fmt.Println(file)
		kind := http.DetectContentType(buff)
		fmt.Println(kind) // do something based on your detection.
		if !stringInSlice(kind, allowedKinds) {
			fmt.Fprintf(w, "Got ya, ivan!\n")
			return
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create(filepath.Join("temp-images/", files[i].Filename))
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func setupRoutes() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadFileHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	setupRoutes()
	fmt.Println("running on :8080")

}

func renderTemplate(w http.ResponseWriter, tmpl string, i interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl+".gohtml", i); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
