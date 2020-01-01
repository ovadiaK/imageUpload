package handler

import (
	"fmt"
	"github.com/ovadiaK/imageUpload/img"
	"net/http"
	"path/filepath"
)

type messageManipulator struct {
	Img   string
	Title string
}

func Manipulator(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		return
	}
	i := r.Form["img"]
	if len(i) != 1 {
		fmt.Println(i)
		return
	}
	name := i[0]
	mess := messageManipulator{}
	mess.Img = filepath.Join("..", img.IMAGE_FOLDER, name)
	mess.Title = name
	renderTemplate(w, "manipulator", mess)

}
